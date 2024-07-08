package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/namnv2496/user-service/internal/wiring"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	app, cleanup, err := wiring.Initilize()
	if err != nil {
		log.Fatalln("Error when init server: ", err)
	}
	defer cleanup()
	if err := app.Start(); err != nil {
		log.Fatalln("Failed to start server")
	}
}
