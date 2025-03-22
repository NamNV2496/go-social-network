package controller

import (
	"context"
	"fmt"
	"net"
	"os"

	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/service"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"google.golang.org/grpc/reflection"
)

type IController interface {
	Start(ctx context.Context) error
}

type Controller struct {
	postv1.UnimplementedPostServiceServer
	likeService    service.ILikeService
	postService    service.IPostService
	commentService service.ICommentService
}

func NewController(
	likeService service.ILikeService,
	postService service.IPostService,
	commentService service.ICommentService,
) IController {
	return &Controller{
		likeService:    likeService,
		postService:    postService,
		commentService: commentService,
	}
}

func (c *Controller) Start(ctx context.Context) error {
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

	var opts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	server := grpc.NewServer(opts...)
	reflection.Register(server)
	postv1.RegisterPostServiceServer(server, c)

	fmt.Printf("gRPC server is running on %s\n", postServiceAddr)
	return server.Serve(listener)
}
