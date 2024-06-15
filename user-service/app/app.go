package app

import (
	"context"

	"github.com/namnv2496/user-service/internal/handler/grpc"
)

type AppInterface interface {
	Start() error
}

type App struct {
	grpcServer grpc.Server
}

func NewApp(
	grpcServer grpc.Server,
) *App {

	return &App{
		grpcServer: grpcServer,
	}
}

func (a App) Start() error {
	return a.grpcServer.Start(context.Background())
}
