package main

import "github.com/namnv2496/user-service/internal/wiring"

func main() {

	app, cleanup, err := wiring.Initilize()
	if err != nil {
		panic("Error when start server")
	}
	defer cleanup()
	app.Start()
}
