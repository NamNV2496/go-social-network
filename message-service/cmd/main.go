package main

import (
	"log"

	"github.com/namnv2496/message-service/internal/handler/route"
	wsHandler "github.com/namnv2496/message-service/internal/handler/ws"
)

func main() {

	hub := wsHandler.NewHub()
	wsHandle := wsHandler.NewHandler(hub)
	go hub.Run()
	route.InitRoute(wsHandle)
	if err := route.Start(":8081"); err != nil {
		log.Fatalln("Failed to start server with port: 8081")
	}
}
