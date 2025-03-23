package cmd

import (
	"context"
	// "fmt"
	"log"
	// "net"
	// "net/http"
	"os"
	"time"

	// "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/controller"
	"github.com/namnv2496/user-service/internal/repository/cache"
	"github.com/namnv2496/user-service/internal/repository/database"
	"github.com/namnv2496/user-service/internal/repository/elasticsearch"
	"github.com/namnv2496/user-service/internal/repository/gateway"
	"github.com/namnv2496/user-service/internal/repository/repo"
	"github.com/namnv2496/user-service/internal/service"
	pb "github.com/namnv2496/user-service/pkg/user_core/v1"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "run service",
	Long:  "run service",
	Run: func(cmd *cobra.Command, args []string) {
		Invoke(
			// StartGRPC,
			// StartREST,
			startServer,
		).Run()
	},
}

func Invoke(invokers ...any) *fx.App {
	conf, _ := configs.NewConfig()
	log.Printf("[config]: %+v", conf)
	app := fx.New(
		fx.StartTimeout(time.Second*10),
		fx.StopTimeout(time.Second*10),
		fx.Provide(
			fx.Annotate(controller.NewGrpcHander, fx.As(new(pb.AccountServiceServer))),
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

func startServer(
	lc fx.Lifecycle,
	grpcClient pb.AccountServiceServer,
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
		SetHTTPEnable(true).
		SetGRPCRegisterFunc(func(s *grpc.Server) {
			pb.RegisterAccountServiceServer(s, grpcClient)
		}).
		SetHTTPRegisterFunc(func(mux *runtime.ServeMux, conn *grpc.ClientConn) {
			pb.RegisterAccountServiceHandler(context.Background(), mux, conn)
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

// func StartGRPC(
// 	lc fx.Lifecycle,
// 	grpcClient pb.AccountServiceServer,
// ) {
// 	lis, err := net.Listen("tcp", ":9090")
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
// 	pb.RegisterAccountServiceServer(grpcServer, grpcClient)
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
// 				err = pb.RegisterAccountServiceHandler(context.Background(), mux, conn)
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
