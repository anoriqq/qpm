package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/anoriqq/qpm"
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
	var aquiferPath string

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

			if aquiferPath != "" {
				c.AquiferPath = aquiferPath
			}

			aquiferPath := os.ExpandEnv(c.AquiferPath)

			oldDir := fmt.Sprintf("%s.old_%s", aquiferPath, time.Now().Format("20060102150405"))
			if err := os.Rename(aquiferPath, oldDir); err != nil && !os.IsNotExist(err) {
				return err
			}

			cl, err := git.NewClient(c.GitHubUsername, c.GitHubToken)
			if err != nil {
				return err
			}

			if err = cl.Clone(aquiferPath, c.AquiferRemote.String()); err != nil {
				return err
			}

			return nil
		},
	}

	aquiferValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "validate specific stratum of aquifer",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			path, err := config.InitConfigFile()
			if err != nil {
				return err
			}

			c, err := config.ReadConfig(path)
			if err != nil {
				return err
			}

			if aquiferPath != "" {
				c.AquiferPath = aquiferPath
			}

			if _, err := qpm.ReadStratum(c, args[0]); err != nil {
				return err
			}

			return nil
		},
	}

	aquiferCmd.PersistentFlags().StringVarP(&aquiferPath, "aquifer", "a", "", "Aquifer directory path")
	aquiferCmd.AddCommand(aquiferPullCmd)
	aquiferCmd.AddCommand(aquiferValidateCmd)
}
