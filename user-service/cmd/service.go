package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/controller"
	"github.com/namnv2496/user-service/internal/repository/cache"
	"github.com/namnv2496/user-service/internal/repository/database/migrations"
	"github.com/namnv2496/user-service/internal/repository/elasticsearch"
	"github.com/namnv2496/user-service/internal/repository/gateway"
	"github.com/namnv2496/user-service/internal/repository/repo"
	"github.com/namnv2496/user-service/internal/repository/sms"
	"github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func Invoke(invokers ...any) *fx.App {
	conf, _ := configs.NewConfig()
	log.Printf("[config]: %+v", conf)
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(
			fx.Annotate(controller.NewGrpcHander, fx.As(new(userv1.AccountServiceServer))),
			fx.Annotate(controller.NewEmailTemplateHander, fx.As(new(userv1.EmailTemplateServiceServer))),
			fx.Annotate(cache.NewRedisClient, fx.As(new(cache.Client))),
			migrations.NewDatabase,
			migrations.InitializeGoquDB,
			fx.Annotate(elasticsearch.NewElasticSearch, fx.As(new(elasticsearch.ElasticSearchClient))),
			fx.Annotate(service.NewUserService, fx.As(new(service.UserService))),
			fx.Annotate(repo.NewUserRepository, fx.As(new(repo.UserRepo))),
			fx.Annotate(repo.NewUserUserRepository, fx.As(new(repo.UserUserRepo))),
			fx.Annotate(service.NewEmailClient, fx.As(new(service.IEmail))),
			fx.Annotate(service.NewOTPService, fx.As(new(service.IOTP))),
			fx.Annotate(sms.NewSms, fx.As(new(sms.ISms))),
			fx.Annotate(repo.NewEmailTemplateRepository, fx.As(new(repo.IEmailTemplateRepo))),
		),
		fx.Supply(
			conf,
		),
		fx.Invoke(invokers...),
	)
	return app
}

// =========================== WAY 1 ===========================
func startServer(
	lc fx.Lifecycle,
	accountServer userv1.AccountServiceServer,
	emailServer userv1.EmailTemplateServiceServer,
) {
	var userServiceAddr string
	if value := os.Getenv("USER_URL"); value != "" {
		userServiceAddr = value
	} else {
		userServiceAddr = "0.0.0.0:5610"
	}
	fmt.Printf("start with port: %s\n", userServiceAddr)
	config := gateway.NewServerConfig().
		SetGRPCAddress(userServiceAddr).
		SetHTTPAddress("localhost:9089").
		SetGRPCEnable(true).
		SetHTTPEnable(true).
		SetGRPCRegisterFunc(func(server *grpc.Server) {
			userv1.RegisterAccountServiceServer(server, accountServer)
			userv1.RegisterEmailTemplateServiceServer(server, emailServer)
		}).
		SetHTTPRegisterFunc(func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
			userv1.RegisterAccountServiceHandler(context.Background(), mux, conn)
			userv1.RegisterEmailTemplateServiceHandler(context.Background(), mux, conn)
		})
	server, err := gateway.NewServer(config)
	if err != nil {
		panic(err)
	}
	// metric.InitPrometheus()
	// expose for prometheus. it will trigger every 15s config in prometheus.yaml
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9999", nil) // expose port 9090
	}()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := server.Serve(ctx)
			if err != nil {
				panic(err)
			}
			return nil
		},
		OnStop: server.Stop,
	})
}

// =========================== WAY 2 ===========================
// func StartGRPC(
// 	lc fx.Lifecycle,
// 	grpcImpl userv1.AccountServiceServer,
// ) {
// 	lis, err := net.Listen("tcp", ":5610")
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	var grpcOpts = []grpc.ServerOption{
// 		grpc.ChainUnaryInterceptor(
// 			validator.UnaryServerInterceptor(),
// 		),
// 		grpc.ChainStreamInterceptor(
// 			validator.StreamServerInterceptor(),
// 		),
// 	}
// 	grpcServer := grpc.NewServer(grpcOpts...)
// 	userv1.RegisterAccountServiceServer(grpcServer, grpcImpl)
// 	reflection.Register(grpcServer)

// 	log.Printf("GRPC server is running on %s", lis.Addr().String())
// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			go func() {
// 				if err := grpcServer.Serve(lis); err != nil {
// 					log.Fatalf("failed to serve gRPC: %v", err)
// 				}
// 			}()
// 			return nil
// 		},
// 		OnStop: func(ctx context.Context) error {
// 			grpcServer.GracefulStop()
// 			return nil
// 		},
// 	})
// }

// func StartREST(
// 	lc fx.Lifecycle,
// ) {
// 	lc.Append(fx.Hook{
// 		OnStart: func(ctx context.Context) error {
// 			go func() {
// 				mux := runtime.NewServeMux()

// 				// Establish a connection to the gRPC server
// 				conn, err := grpc.NewClient(":5610", grpc.WithInsecure())
// 				if err != nil {
// 					log.Fatalf("failed to connect to gRPC server: %v", err)
// 				}
// 				defer conn.Close()

// 				// Register the gRPC service with the HTTP gateway
// 				err = userv1.RegisterAccountServiceHandler(context.Background(), mux, conn)
// 				if err != nil {
// 					log.Fatalf("failed to register gRPC service with HTTP gateway: %v", err)
// 				}

// 				// Start the HTTP server
// 				const port = "8083"
// 				fmt.Printf("HTTP server is running on %s\n", port)
// 				log.Fatal(http.ListenAndServe(":"+port, mux))
// 			}()
// 			return nil
// 		},
// 		OnStop: func(ctx context.Context) error {
// 			return nil
// 		},
// 	})
// }
