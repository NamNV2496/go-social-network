//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/post-service/app"
	"github.com/namnv2496/post-service/internal/configs"
	"github.com/namnv2496/post-service/internal/controller"
	"github.com/namnv2496/post-service/internal/repository"
	"github.com/namnv2496/post-service/internal/service"
)

func Initilize() (*app.App, error) {
	wire.Build(
		configs.ConfigWireSet,
		repository.RepositoryWireSet,
		service.ServiceWireSet,
		controller.ControllerWireSet,
		app.NewApp,
	)
	return nil, nil
}
