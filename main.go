package main

import (
	"os"

	"github.com/anoriqq/qpm/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
