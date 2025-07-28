package gateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var httpRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total-test",
		Help: "Total number of HTTP requests.",
	},
	[]string{"path", "method"},
)

type httpServerConfig struct {
	addr         string
	registerFunc func(mux *runtime.ServeMux, conn *grpc.ClientConn)
	middlewares  []func(http.Handler) http.Handler
	enable       bool
}

type httpServer struct {
	servers *http.Server
	conn    *grpc.ClientConn
}

func newHTTPServer(conf *config) (*httpServer, error) {
	prometheus.MustRegister(httpRequestsTotal)
	if !conf.http.enable {
		return nil, nil
	}
	if conf.grpc.addr == "" {
		return nil, fmt.Errorf("register http Handler is required")
	}
	options := []runtime.ServeMuxOption{}
	mux := runtime.NewServeMux(options...)

	handler := http.Handler(mux)
	conn, err := grpc.NewClient(conf.grpc.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	conf.http.registerFunc(mux, conn)
	fmt.Println("start http server: ", conf.http.addr)
	return &httpServer{
		servers: &http.Server{
			Addr:    conf.http.addr,
			Handler: handler,
		},
		conn: conn,
	}, nil
}

func (s *httpServer) Serve() error {
	if err := s.servers.ListenAndServe(); err != nil {
		return fmt.Errorf("http start failed: %w", err)
	}
	return nil
}

func (s *httpServer) Stop(ctx context.Context) {
	defer s.conn.Close()
	if err := s.servers.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server Shutdown: %v", err)
	}
}
