package cmd

import (
	"github.com/anoriqq/qpm/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qpm",
	Short: "qpm",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		err := config.InitConfig()
		if err != nil {
			panic(err)
		}
	})
}
