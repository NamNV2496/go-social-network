package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/namnv2496/newsfeed-service/internal/handler/consumers"
	"github.com/namnv2496/newsfeed-service/internal/handler/grpc"
)

type AppInterface interface {
	Start() error
}

type App struct {
	grpcServer grpc.Server
	consumer   consumers.ConsumerHandler
}

func NewApp(
	grpcServer grpc.Server,
	consumer consumers.ConsumerHandler,
) *App {

	return &App{
		grpcServer: grpcServer,
		consumer:   consumer,
	}
}

func (a App) Start() error {

	go func() {
		if err := a.consumer.StartConsumerUp(context.Background()); err != nil {
			log.Fatalln("Failed to start consumer")
		}
	}()

	go func() {
		if err := a.grpcServer.Start(context.Background()); err != nil {
			log.Fatalln("Failed to start grpc server")
		}
	}()
	BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	return nil
}

func BlockUntilSignal(signals ...os.Signal) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, signals...)
	<-done
}
