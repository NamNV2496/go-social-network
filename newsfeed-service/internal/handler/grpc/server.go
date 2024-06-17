package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	newsfeedv1 "github.com/namnv2496/newsfeed-service/internal/handler/generated/newsfeed_core/v1"

	"google.golang.org/grpc"
)

type Server interface {
	Start(ctx context.Context) error
}

type server struct {
	handler newsfeedv1.NewsfeedServiceServer
}

func NewServer(
	handler newsfeedv1.NewsfeedServiceServer,
) Server {
	return &server{
		handler: handler,
	}
}

func (s server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", "localhost:5602")
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
	newsfeedv1.RegisterNewsfeedServiceServer(server, s.handler)

	fmt.Printf("gRPC server is running on %s\n", "localhost:5602")
	return server.Serve(listener)
}
