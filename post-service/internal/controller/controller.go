package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"net"
	"os"

	"github.com/namnv2496/post-service/configs"
	postv1 "github.com/namnv2496/post-service/internal/handler/generated/post_core/v1"
	"github.com/namnv2496/post-service/internal/pkg/logger"
	"github.com/namnv2496/post-service/internal/pkg/metric"
	"github.com/namnv2496/post-service/internal/service"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc/reflection"
)

type IController interface {
	Start(ctx context.Context) error
}

type Controller struct {
	postv1.PostServiceServer
	logger         *logger.Logger
	likeService    service.ILikeService
	postService    service.IPostService
	commentService service.ICommentService
}

func NewController(
	cfg configs.Config,
	likeService service.ILikeService,
	postService service.IPostService,
	commentService service.ICommentService,
) IController {
	return &Controller{
		logger:         logger.NewLogger("post-controller", "post"),
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
	postv1.RegisterPostServiceServer(server, c)

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
