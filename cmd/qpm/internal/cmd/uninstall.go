package cmd

import (
	"bufio"
	"os"

	"github.com/anoriqq/qpm"
	"github.com/anoriqq/qpm/cmd/qpm/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	var aquiferPath string
	var shell string

	uninstallCmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall specific package",
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

			cfg := c
			cfg.Shell = []string{shell}

			return qpm.Execute(cfg, s, qpm.Uninstall, bufio.NewWriter(os.Stdout), bufio.NewWriter(os.Stderr))
		},
	}

	uninstallCmd.PersistentFlags().StringVarP(&aquiferPath, "aquifer", "a", "", "Aquifer directory path")
	uninstallCmd.PersistentFlags().StringVarP(&shell, "shell", "s", "", "Shell")
	rootCmd.AddCommand(uninstallCmd)
}
