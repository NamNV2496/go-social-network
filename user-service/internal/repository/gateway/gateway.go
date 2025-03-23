package gateway

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/namnv2496/user-service/internal/repository/logger"
	"google.golang.org/grpc"
)

type config struct {
	log  logger.ILogger
	grpc grpcServerConfig
	http httpServerConfig
}

func NewServerConfig() *config {
	conf := &config{
		grpc: grpcServerConfig{
			addr:         "localhost:9090",
			registerFunc: nil,
			enable:       false,
		},
		http: httpServerConfig{
			addr:         "localhost:8080",
			registerFunc: nil,
			middlewares:  nil,
			enable:       false,
		},
		log: logger.NewLogger("gateway", ""),
	}
	return conf
}

func (c *config) SetGRPCAddress(addr string) *config {
	c.grpc.addr = addr
	return c
}

func (c *config) SetGRPCRegisterFunc(registerFunc func(s *grpc.Server)) *config {
	c.grpc.registerFunc = registerFunc
	return c
}

func (c *config) SetGRPCEnable(grpcEnable bool) *config {
	c.grpc.enable = grpcEnable
	return c
}

func (c *config) SetHTTPAddress(addr string) *config {
	c.http.addr = addr
	return c
}

func (c *config) SetHTTPRegisterFunc(registerFunc func(mux *runtime.ServeMux, conn *grpc.ClientConn)) *config {
	c.http.registerFunc = registerFunc
	return c
}

func (c *config) SetHTTPEnable(httpEnable bool) *config {
	c.http.enable = httpEnable
	return c
}

func (c *config) SetMiddleware(middlerwares []func(http.Handler) http.Handler) *config {
	c.http.middlewares = middlerwares
	return c
}
