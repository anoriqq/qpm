package cmd

import (
	"github.com/anoriqq/qpm/internal/service/getscript"
	"github.com/spf13/cobra"
)

var getscriptCmd = &cobra.Command{
	Use:   "getscript",
	Short: "get script dir form remote repository",
	RunE:  getscriptRun,
}

func init() {
	rootCmd.AddCommand(getscriptCmd)
}

func getscriptRun(_ *cobra.Command, _ []string) error {
	return getscript.GetScript()
}
