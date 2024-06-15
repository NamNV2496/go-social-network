package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/post-service/internal/handler/grpc"
)

var HandlerWireSet = wire.NewSet(
	grpc.NewGrpcHander,
	grpc.NewServer,
)
