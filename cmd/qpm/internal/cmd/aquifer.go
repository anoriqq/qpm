package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/anoriqq/qpm/cmd/qpm/internal/config"
	"github.com/anoriqq/qpm/cmd/qpm/internal/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aquiferCmd)
}

var aquiferCmd = &cobra.Command{
	Use:   "aquifer",
	Short: "Manage aquifer",
}

func init() {
	aquiferPullCmd := &cobra.Command{
		Use:   "pull",
		Short: "Get aquifer form remote repository",
		RunE: func(_ *cobra.Command, _ []string) error {
			path, err := config.InitConfigFile()
			if err != nil {
				return err
			}

			c, err := config.ReadConfig(path)
			if err != nil {
				return err
			}

			oldDir := fmt.Sprintf("%s.old_%s", c.AquiferPath, time.Now().Format("20060102150405"))
			if err := os.Rename(c.AquiferPath, oldDir); err != nil && !os.IsNotExist(err) {
				return err
			}

			cl, err := git.NewClient(c.GitHubUsername, c.GitHubToken)
			if err != nil {
				return err
			}

			if err = cl.Clone(c.AquiferPath, c.AquiferRemote.String()); err != nil {
				return err
			}

			return nil
		},
	}

	aquiferCmd.AddCommand(aquiferPullCmd)
}
