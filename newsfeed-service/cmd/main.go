package main

import (
	"log"

	"github.com/namnv2496/newsfeed-service/internal/pkg/metric"
	"github.com/namnv2496/newsfeed-service/internal/wiring"
)

func main() {
	metric.InitPrometheus()
	app, cleanup, err := wiring.Initilize()
	if err != nil {
		log.Fatalln("Error when start server: ", err)
	}
	defer cleanup()
	if err := app.Start(); err != nil {
		log.Fatalln("Failed to start server")
	}
}
