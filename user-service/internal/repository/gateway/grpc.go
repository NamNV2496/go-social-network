package gateway

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
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
	}
	var grpcOpts = []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			validator.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			validator.StreamServerInterceptor(),
		),
	}
	server := grpc.NewServer(grpcOpts...)
	conf.grpc.registerFunc(server)
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
