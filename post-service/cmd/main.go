package main

import (
	"log"

	"github.com/namnv2496/post-service/internal/wiring"
)

func main() {

	app, cleanup, err := wiring.Initilize()
	if err != nil {
		log.Fatalln("Error when init server: ", err)
	}
	defer cleanup()
	if err := app.Start(); err != nil {
		log.Fatalln("Failed to start server: ", err)
	}
}
