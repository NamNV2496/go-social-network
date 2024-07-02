package http

import (
	"context"
	"log"
	"net/http"
	"os"

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

	var userServiceAddr string
	if value := os.Getenv("USER_URL"); value != "" {
		userServiceAddr = value
	} else {
		userServiceAddr = "0.0.0.0:5600"
	}

	conn, err := grpc.NewClient(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to order service: %v", err)
	}
	defer conn.Close()
	return nil
}

func (s *server) ConnectToPostService(ctx context.Context) error {

	var postServiceAddr string
	if value := os.Getenv("POST_URL"); value != "" {
		postServiceAddr = value
	} else {
		postServiceAddr = "0.0.0.0:5601"
	}
	conn, err := grpc.NewClient(postServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	// TODO
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

	var userServiceAddr string
	if value := os.Getenv("USER_URL"); value != "" {
		userServiceAddr = value
	} else {
		userServiceAddr = s.grpcConfig.UserServiceAddress
	}

	var postServiceAddr string
	if value := os.Getenv("POST_URL"); value != "" {
		postServiceAddr = value
	} else {
		postServiceAddr = s.grpcConfig.PostServiceAddress
	}

	var newsfeedServiceAddr string
	if value := os.Getenv("NEWSFEED_URL"); value != "" {
		newsfeedServiceAddr = value
	} else {
		newsfeedServiceAddr = s.grpcConfig.NewfeedsServiceAddress
	}

	log.Println("Connect to user: ", userServiceAddr)
	log.Println("Connect to post: ", postServiceAddr)
	log.Println("Connect to newsfeed: ", newsfeedServiceAddr)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	var err error

	// connect to user-service
	err = userv1.RegisterAccountServiceHandlerFromEndpoint(
		ctx,
		mux,
		userServiceAddr,
		opts,
	)
	if err != nil {
		return err
	}

	// connect to post-service
	err = postv1.RegisterPostServiceHandlerFromEndpoint(
		ctx,
		mux,
		postServiceAddr,
		opts,
	)
	if err != nil {
		return err
	}
	// connect to newsfeed-service
	err = newsfeedv1.RegisterNewsfeedServiceHandlerFromEndpoint(
		ctx,
		mux,
		newsfeedServiceAddr,
		opts,
	)
	if err != nil {
		log.Println("Error: ", err)
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
