package main

import (
	"github.com/kelveny/mockcompose/cmd"
)

func main() {
	cmd.Log(cmd.PROMPT, "Running mockcompose version 0.1.4\n")

	cmd.Execute()
}
