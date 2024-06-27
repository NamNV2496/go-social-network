package main

import (
	"log"

	"github.com/namnv2496/newsfeed-service/internal/wiring"
)

func main() {

	app, cleanup, err := wiring.Initilize()
	if err != nil {
		panic("Error when start server")
	}
	defer cleanup()
	if err := app.Start(); err != nil {
		log.Fatalln("Failed to start server")
	}
}
