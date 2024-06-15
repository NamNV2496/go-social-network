package main

import (
	"context"
	"fmt"

	"github.com/namnv2496/http_gateway/internal/wiring"
)

func main() {

	server, err := wiring.Initialize()
	if err != nil {
		panic("Cannot init server")
	}
	fmt.Println("Server is running...")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server.ConnectToUserService(ctx)
	server.Start("rest")
}
