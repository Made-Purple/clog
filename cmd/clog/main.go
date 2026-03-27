package main

import (
	"os"

	"github.com/made-purple/clog/internal/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
