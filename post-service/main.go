package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/controller"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg"
	"github.com/namnv2496/post-service/internal/pkg/metric"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/repository/database"
	"github.com/namnv2496/post-service/internal/repository/mq/producer"
	"github.com/namnv2496/post-service/internal/service"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := execute()
	if err != nil {
		panic(err)
	}
}

func execute() error {
	var root = &cobra.Command{
		Short: "root",
		Run: func(cmd *cobra.Command, args []string) {
			Invoke(
				startServer,
			).Run()
		},
	}
	// separate command if needed
	// root.AddCommand(command.NotificationCmd)
	return root.Execute()
}

func Invoke(invokers ...any) *fx.App {
	conf, _ := configs.NewConfig()
	log.Printf("[config]: %+v", conf)
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(
			fx.Annotate(controller.NewController, fx.As(new(postv1.PostServiceServer))),
			// fx.Annotate(database.NewDatabaseConnection, fx.As(new(*gorm.DB))),
			database.NewDatabaseConnection,
			fx.Annotate(producer.NewClient, fx.As(new(producer.Client))),

			fx.Annotate(repository.NewLikeRepository, fx.As(new(repository.ILikeRepository))),
			fx.Annotate(repository.NewLikeCountRepository, fx.As(new(repository.ILikeCountRepository))),
			fx.Annotate(repository.NewCommentRepository, fx.As(new(repository.ICommentRepository))),
			fx.Annotate(repository.NewPostRepository, fx.As(new(repository.IPostRepository))),
			fx.Annotate(repository.NewTransaction, fx.As(new(repository.ITransaction))),
			fx.Annotate(repository.NewCommentRuleRepository, fx.As(new(repository.ICommentRuleRepository))),

			fx.Annotate(service.NewPostService, fx.As(new(service.IPostService))),
			fx.Annotate(service.NewCommentService, fx.As(new(service.ICommentService))),
			fx.Annotate(service.NewLikeService, fx.As(new(service.ILikeService))),

			fx.Annotate(pkg.NewTrie, fx.As(new(pkg.ITrie))),
			// notification
			fx.Annotate(controller.NewNotificationController, fx.As(new(postv1.NotificationServiceServer))),
			fx.Annotate(repository.NewNotificationRepository, fx.As(new(repository.INotificationRepository))),
			fx.Annotate(service.NewNotificationService, fx.As(new(service.INotificationService))),
		),
		fx.Supply(
			conf,
		),
		fx.Invoke(invokers...),
	)
	return app
}

func startServer(
	lc fx.Lifecycle,
	controllerClient postv1.PostServiceServer,
	notificationController postv1.NotificationServiceServer,
) error {
	var postServiceAddr string
	if value := os.Getenv("POST_URL"); value != "" {
		postServiceAddr = value
	} else {
		postServiceAddr = "localhost:5611"
	}

	listener, err := net.Listen("tcp", postServiceAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	metric.InitPrometheus()
	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			validator.UnaryServerInterceptor(),
			// interceptors.MetricsInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
		),
	}
	server := grpc.NewServer(opts...)
	postv1.RegisterPostServiceServer(server, controllerClient)
	postv1.RegisterNotificationServiceServer(server, notificationController)

	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// register prometheus
	grpc_prometheus.EnableHandlingTimeHistogram()
	// Register Prometheus metrics handler.
	grpc_prometheus.Register(server)
	// reflection
	reflection.Register(server)

	fmt.Printf("gRPC server is running on %s\n", postServiceAddr)

	// expose promethus for debugging
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9998", nil))
	}()

	return server.Serve(listener)
}
