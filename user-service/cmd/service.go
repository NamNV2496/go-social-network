package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/controller"
	"github.com/namnv2496/user-service/internal/repository/cache"
	"github.com/namnv2496/user-service/internal/repository/database"
	"github.com/namnv2496/user-service/internal/repository/elasticsearch"
	"github.com/namnv2496/user-service/internal/repository/gateway"
	"github.com/namnv2496/user-service/internal/repository/repo"
	"github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Invoke(invokers ...any) *fx.App {
	conf, _ := configs.NewConfig()
	log.Printf("[config]: %+v", conf)
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(
			fx.Annotate(controller.NewGrpcHander, fx.As(new(userv1.AccountServiceServer))),
			fx.Annotate(cache.NewRedisClient, fx.As(new(cache.Client))),
			database.NewDatabase,
			database.InitializeGoquDB,
			fx.Annotate(elasticsearch.NewElasticSearch, fx.As(new(elasticsearch.ElasticSearchClient))),
			fx.Annotate(service.NewUserService, fx.As(new(service.UserService))),
			fx.Annotate(repo.NewUserRepository, fx.As(new(repo.UserRepo))),
			fx.Annotate(repo.NewUserUserRepository, fx.As(new(repo.UserUserRepo))),
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
	grpcServer userv1.AccountServiceServer,
) {
	var userServiceAddr string
	if value := os.Getenv("USER_URL"); value != "" {
		userServiceAddr = value
	} else {
		userServiceAddr = "localhost:9090"
	}
	config := gateway.NewServerConfig().
		SetGRPCAddress("localhost:5610").
		SetHTTPAddress(userServiceAddr).
		SetGRPCEnable(true).
		SetHTTPEnable(false).
		SetGRPCRegisterFunc(func(s *grpc.Server) {
			userv1.RegisterAccountServiceServer(s, grpcServer)
		}).
		SetHTTPRegisterFunc(func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
			userv1.RegisterAccountServiceHandler(context.Background(), mux, conn)
		})
	server, err := gateway.NewServer(config)
	if err != nil {
		panic(err)
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Serve(ctx)
				if err != nil {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: server.Stop,
	})
}

// =========================== WAY 2 ===========================
func StartGRPC(
	lc fx.Lifecycle,
	grpcClient userv1.AccountServiceServer,
) {
	lis, err := net.Listen("tcp", ":5610")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	grpcServer := grpc.NewServer(grpcOpts...)
	userv1.RegisterAccountServiceServer(grpcServer, grpcClient)
	reflection.Register(grpcServer)

	log.Printf("GRPC server is running on %s", lis.Addr().String())
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := grpcServer.Serve(lis); err != nil {
					log.Fatalf("failed to serve gRPC: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
}

func StartREST(
	lc fx.Lifecycle,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				mux := runtime.NewServeMux()

				// Establish a connection to the gRPC server
				conn, err := grpc.NewClient(":5610", grpc.WithInsecure())
				if err != nil {
					log.Fatalf("failed to connect to gRPC server: %v", err)
				}
				defer conn.Close()

				// Register the gRPC service with the HTTP gateway
				err = userv1.RegisterAccountServiceHandler(context.Background(), mux, conn)
				if err != nil {
					log.Fatalf("failed to register gRPC service with HTTP gateway: %v", err)
				}

				// Start the HTTP server
				const port = "8083"
				fmt.Printf("HTTP server is running on %s\n", port)
				log.Fatal(http.ListenAndServe(":"+port, mux))
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
