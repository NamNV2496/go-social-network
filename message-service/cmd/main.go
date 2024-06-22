package main

import (
	"github.com/namnv2496/message-service/internal/handler/route"
	wsHandler "github.com/namnv2496/message-service/internal/handler/ws"
	"github.com/namnv2496/message-service/internal/logic"
)

func main() {

	hub := logic.NewHub()
	wsHandle := wsHandler.NewHandler(hub)
	go hub.Run()
	route.InitRoute(wsHandle)
	route.Start(":8081")
}
