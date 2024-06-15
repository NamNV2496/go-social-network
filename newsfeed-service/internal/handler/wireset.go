package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/newsfeed-service/internal/handler/consumers"
	"github.com/namnv2496/newsfeed-service/internal/handler/grpc"
)

var HandlerWireSet = wire.NewSet(
	grpc.NewGrpcHander,
	grpc.NewServer,
	consumers.NewKafkaHandler,
)
