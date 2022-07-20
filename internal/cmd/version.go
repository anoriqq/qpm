package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	RunE:  versionRun,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

const version = "v0.0.6"

func versionRun(_ *cobra.Command, _ []string) error {
	fmt.Println(version)
	return nil
}
