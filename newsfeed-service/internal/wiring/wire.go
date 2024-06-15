//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/newsfeed-service/app"
	"github.com/namnv2496/newsfeed-service/internal/cache"
	"github.com/namnv2496/newsfeed-service/internal/configs"
	"github.com/namnv2496/newsfeed-service/internal/database"
	"github.com/namnv2496/newsfeed-service/internal/handler"
	"github.com/namnv2496/newsfeed-service/internal/logic"
	"github.com/namnv2496/newsfeed-service/internal/mq"
)

func Initilize() (*app.App, func(), error) {
	wire.Build(
		configs.ConfigWireSet,
		database.DatabaseWireSet,
		logic.LogicWireSet,
		handler.HandlerWireSet,
		mq.MQWireSet,
		cache.CacheWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
