package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install packages",
	RunE:  installRun,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installRun(_ *cobra.Command, _ []string) error {
	return nil
}
