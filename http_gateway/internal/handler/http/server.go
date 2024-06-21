package http

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/namnv2496/http_gateway/internal/configs"
	newsfeedv1 "github.com/namnv2496/http_gateway/internal/handler/generated/newsfeed_core/v1"
	postv1 "github.com/namnv2496/http_gateway/internal/handler/generated/post_core/v1"
	userv1 "github.com/namnv2496/http_gateway/internal/handler/generated/user_core/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	Start(string) error
	ConnectToUserService(ctx context.Context) error
	ConnectToPostService(ctx context.Context) error
}

type server struct {
	gatewayConfig configs.Gateway
	grpcConfig    configs.GRPC
	authConfig    configs.Auth
}

func NewServer(
	httpConfig configs.Gateway,
	grpcConfig configs.GRPC,
	authConfig configs.Auth,
) Server {

	return &server{
		gatewayConfig: httpConfig,
		grpcConfig:    grpcConfig,
		authConfig:    authConfig,
	}
}

func (s *server) ConnectToUserService(ctx context.Context) error {

	orderServiceAddr := "0.0.0.0:5600"
	conn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to order service: %v", err)
	}
	defer conn.Close()
	return nil
}

func (s *server) ConnectToPostService(ctx context.Context) error {

	orderServiceAddr := "0.0.0.0:5601"
	conn, err := grpc.Dial(orderServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to order service: %v", err)
	}
	defer conn.Close()
	return nil
}

func (s *server) Start(serverType string) error {

	if serverType == "grpc" {
		return s.runGRPCServer()
	} else {
		return s.runRESTServer()
	}
}

func (s *server) runGRPCServer() error {
	// interceptor := service.NewAuthInterceptor(jwtManager)
	// serverOptions := []grpc.ServerOption{
	// 	grpc.UnaryInterceptor(interceptor.Unary()),
	// 	grpc.StreamInterceptor(interceptor.Stream()),
	// }

	// server := grpc.NewServer(
	// 	grpc.ChainUnaryInterceptor(
	// 		ErrorLogger(),
	// 		EnsureAuthIsValid,
	// 	),
	// )
	// if enableTLS {
	// 	tlsCredentials, err := loadTLSCredentials()
	// 	if err != nil {
	// 		return fmt.Errorf("cannot load TLS credentials: %w", err)
	// 	}

	// 	serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	// }

	// grpcServer := grpc.NewServer(serverOptions...)

	// log.Printf("Start GRPC server at %s, TLS = %t", s.gatewayConfig.GatewayAddress, )
	// return grpcServer.Serve(s.gatewayConfig.GatewayAddress)
	return nil
}

func (s *server) runRESTServer() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		forwardHeaderToClient(),
		withErrorHandler(),
		withMetadata(),
	)
	// mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	var err error

	// connect to user-service
	err = userv1.RegisterAccountServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.UserServiceAddress,
		opts,
	)

	if err != nil {
		return err
	}

	// connect to post-service
	err = postv1.RegisterPostServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.PostServiceAddress,
		opts,
	)

	// connect to newsfeed-service
	err = newsfeedv1.RegisterNewsfeedServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.grpcConfig.NewfeedsServiceAddress,
		opts,
	)

	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	// start server of http
	return http.ListenAndServe(s.gatewayConfig.GatewayAddress, corsHandler(mux))
}

// custom CORS handler
func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins for simplicity, but adjust for production
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
