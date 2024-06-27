package main

import (
	"log"

	"github.com/namnv2496/http_gateway/internal/wiring"
)

func main() {

	server, err := wiring.Initialize()
	if err != nil {
		log.Fatalln("Cannot init server")
	}
	log.Println("Server is running...")

	if err := server.Start("rest"); err != nil {
		log.Fatalln("Start server REST failed.")
	}
}
