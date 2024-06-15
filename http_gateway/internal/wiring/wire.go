//go:build wireinject

package wiring

import (
	"github.com/google/wire"
	"github.com/namnv2496/http_gateway/internal/configs"
	"github.com/namnv2496/http_gateway/internal/handler"
	"github.com/namnv2496/http_gateway/internal/handler/http"
)

func Initialize() (http.Server, error) {

	wire.Build(
		configs.ConfigWireSet,
		handler.HandlerWireSet,
	)
	return nil, nil
}
