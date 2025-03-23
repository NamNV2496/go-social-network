package main

import "github.com/namnv2496/user-service/cmd"

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
