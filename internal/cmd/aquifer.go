package cmd

import (
	"github.com/anoriqq/qpm/internal/service/aquifer"
	"github.com/spf13/cobra"
)

var aquiferCmd = &cobra.Command{
	Use:   "aquifer",
	Short: "manage aquifer",
}

func init() {
	rootCmd.AddCommand(aquiferCmd)
}

var aquiferPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "get aquifer form remote repository",
	RunE:  aquiferPullRun,
}

func init() {
	aquiferCmd.AddCommand(aquiferPullCmd)
}

func aquiferPullRun(_ *cobra.Command, _ []string) error {
	return aquifer.Pull()
}
