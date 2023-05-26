package cmd

import (
	"github.com/anoriqq/qpm"
	"github.com/anoriqq/qpm/cmd/qpm/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	var aquiferPath string

	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Unnstall specifc package",
		Example: `  # Uninstall foo package
  qpm uninstall foo`,
		Args: cobra.RangeArgs(1, 2),
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

			s, err := qpm.ReadStratum(c, args[0])
			if err != nil {
				return err
			}

			return qpm.Execute(c, s, qpm.Uninstall)
		},
	}

	uninstallCmd.PersistentFlags().StringVarP(&aquiferPath, "aquifer", "a", "", "Aquifer directory path")
	rootCmd.AddCommand(uninstallCmd)
}
