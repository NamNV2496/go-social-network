//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/user-service/app"
	"github.com/namnv2496/user-service/internal/cache"
	"github.com/namnv2496/user-service/internal/configs"
	"github.com/namnv2496/user-service/internal/database"
	"github.com/namnv2496/user-service/internal/handler"
	"github.com/namnv2496/user-service/internal/logic"
	"github.com/namnv2496/user-service/internal/repo"
)

func Initilize() (*app.App, func(), error) {
	wire.Build(
		configs.ConfigWireSet,
		database.DatabaseWireSet,
		logic.LogicWireSet,
		handler.HandlerWireSet,
		cache.CacheWireSet,
		repo.RepoWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
