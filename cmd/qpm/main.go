package main

import (
	"os"

	"github.com/anoriqq/qpm/cmd/qpm/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
