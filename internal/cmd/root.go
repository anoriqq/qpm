package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qpm",
	Short: "qpm",
}

func Execute() error {
	return rootCmd.Execute()
}
