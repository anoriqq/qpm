package main

import (
	"fmt"
	"os"

	"github.com/anoriqq/qpm/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}
