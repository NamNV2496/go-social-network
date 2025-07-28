package gateway

import (
	"context"
	"fmt"
	"net"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type grpcServerConfig struct {
	addr         string
	registerFunc func(s *grpc.Server)
	enable       bool
}

type grpcServer struct {
	servers *grpc.Server
	lis     net.Listener
}

func newGRPCServer(conf *config) (*grpcServer, error) {
	if !conf.grpc.enable {
		return nil, nil
	}
	if conf.grpc.registerFunc == nil {
		return nil, fmt.Errorf("register GRPC Handler is required")
	}
	if conf.grpc.addr == "" {
		return nil, fmt.Errorf("address GRPC Handler is required")
	}

	wrapped := grpc_prometheus.UnaryServerInterceptor

	debugInterceptor := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		fmt.Println("GRPC CALL:", info.FullMethod)
		return wrapped(ctx, req, info, handler)
	}
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(debugInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)
	conf.grpc.registerFunc(server)

	// register server to prometheus
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	// register prometheus
	grpc_prometheus.EnableHandlingTimeHistogram()
	// Register Prometheus metrics handler.
	grpc_prometheus.Register(server)
	// reflection
	reflection.Register(server)

	lis, err := net.Listen("tcp", conf.grpc.addr)
	if err != nil {
		return nil, err
	}
	fmt.Println("start grpc server: ", conf.grpc.addr)
	return &grpcServer{
		servers: server,
		lis:     lis,
	}, nil
}

func (s *grpcServer) Serve() error {
	if err := s.servers.Serve(s.lis); err != nil {
		return fmt.Errorf("grpc start failed: %w", err)
	}
	return nil
}

func (s *grpcServer) Stop() {
	s.servers.GracefulStop()
}
