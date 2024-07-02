package grpc

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"

	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler postv1.PostServiceServer
}

func NewServer(
	handler postv1.PostServiceServer,
) Server {
	return &server{
		handler: handler,
	}
}

func (s server) Start(ctx context.Context) error {

	var postServiceAddr string
	if value := os.Getenv("POST_URL"); value != "" {
		postServiceAddr = value
	} else {
		postServiceAddr = "localhost:5601"
	}

	listener, err := net.Listen("tcp", postServiceAddr)
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
	postv1.RegisterPostServiceServer(server, s.handler)

	fmt.Printf("gRPC server is running on %s\n", postServiceAddr)
	return server.Serve(listener)
}
