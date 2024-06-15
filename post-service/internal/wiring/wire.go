//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/post-service/app"
	"github.com/namnv2496/post-service/internal/configs"
	"github.com/namnv2496/post-service/internal/database"
	"github.com/namnv2496/post-service/internal/handler"
	"github.com/namnv2496/post-service/internal/logic"
	"github.com/namnv2496/post-service/internal/mq"
)

func Initilize() (*app.App, func(), error) {
	wire.Build(
		configs.ConfigWireSet,
		database.DatabaseWireSet,
		logic.LogicWireSet,
		handler.HandlerWireSet,
		mq.MQWireSet,
		app.NewApp,
	)
	return nil, nil, nil
}
