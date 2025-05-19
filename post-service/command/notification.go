package command

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/namnv2496/post-service/configs"
	"github.com/namnv2496/post-service/internal/controller"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg/metric"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/repository/database"
	"github.com/namnv2496/post-service/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var NotificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "init database",
	Long:  "init database",
	Run: func(cmd *cobra.Command, args []string) {
		InvokeNoti(
			starNoti,
		).Run()
	},
}

func InvokeNoti(invokers ...any) *fx.App {
	conf, _ := configs.NewConfig()
	log.Printf("[config]: %+v", conf)
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(
			fx.Annotate(controller.NewNotificationController, fx.As(new(postv1.NotificationServiceServer))),
			database.NewDatabaseConnection,
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

func starNoti(
	lc fx.Lifecycle,
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
			validator.UnaryServerInterceptor(),
			// interceptors.MetricsInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	server := grpc.NewServer(opts...)
	reflection.Register(server)
	postv1.RegisterNotificationServiceServer(server, notificationController)

	fmt.Printf("gRPC server is running on %s\n", postServiceAddr)

	// expose promethus for debugging
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, World!")
		})
		http.Handle("/metrics", promhttp.Handler())

		log.Fatal(http.ListenAndServe(":8090", nil))
	}()

	return server.Serve(listener)
}
