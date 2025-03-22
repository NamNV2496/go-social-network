package app

import (
	"context"

	"github.com/namnv2496/post-service/internal/controller"
)

type AppInterface interface {
	Start() error
}

type App struct {
	controller controller.IController
}

func NewApp(
	controller controller.IController,
) *App {
	return &App{
		controller: controller,
	}
}

func (a App) Start() error {
	return a.controller.Start(context.Background())
}
