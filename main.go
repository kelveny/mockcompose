package main

import (
	"fmt"

	"github.com/kelveny/mockcompose/cmd"
)

func main() {
	fmt.Println("Running docker compose version 0.1.0-21042602")
	fmt.Println()

	cmd.Execute()
}
