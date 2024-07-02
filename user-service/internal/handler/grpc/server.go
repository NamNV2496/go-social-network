package grpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"

	userv1 "github.com/namnv2496/user-service/internal/handler/generated/user_core/v1"
	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler userv1.AccountServiceServer
}

func NewServer(
	handler userv1.AccountServiceServer,
) Server {
	return &server{
		handler: handler,
	}
}

// Start implements Server.
func (s server) Start(ctx context.Context) error {

	var userServiceAddr string
	if value := os.Getenv("USER_URL"); value != "" {
		userServiceAddr = value
	} else {
		userServiceAddr = "localhost:5600"
	}

	listener, err := net.Listen("tcp", userServiceAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	server := grpc.NewServer(opts...)
	userv1.RegisterAccountServiceServer(server, s.handler)

	fmt.Printf("gRPC server is running on %s\n", userServiceAddr)
	return server.Serve(listener)
}
