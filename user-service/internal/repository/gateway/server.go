package gateway

import (
	"context"
)

type Server struct {
	grpc *grpcServer
	http *httpServer
}

func NewServer(
	conf *config,
) (*Server, error) {
	grpcServer, err := newGRPCServer(conf)
	if err != nil {
		return nil, err
	}
	httpServer, err := newHTTPServer(conf)
	if err != nil {
		return nil, err
	}
	return &Server{
		grpc: grpcServer,
		http: httpServer,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	ch := make(chan error)
	defer close(ch)

	go func() {
		if s.grpc != nil {
			ch <- s.grpc.Serve()
		}
	}()
	go func() {
		if s.http != nil {
			ch <- s.http.Serve()
		}
	}()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	defer s.grpc.Stop()
	defer s.http.Stop(ctx)
	return nil
}
