package cmd

import (
	"github.com/anoriqq/qpm"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "qpm",
	Short:   "Qanat Package Manager",
	Version: qpm.Version(),
}

func Execute() error {
	return rootCmd.Execute()
}
