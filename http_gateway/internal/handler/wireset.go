package handler

import (
	"github.com/google/wire"
	"github.com/namnv2496/http_gateway/internal/handler/http"
)

var HandlerWireSet = wire.NewSet(
	http.NewServer,
)
