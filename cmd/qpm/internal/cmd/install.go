package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/anoriqq/qpm"
	"github.com/anoriqq/qpm/cmd/qpm/internal/config"
	"github.com/anoriqq/qpm/cmd/qpm/internal/survey"
	"github.com/spf13/cobra"
)

func init() {
	var aquiferPath string
	var shell string

	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install specific package",
		Example: `  # Install foo package
  qpm install foo`,
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

			if alreadyInstalled, err := qpm.IsAlreadyInstalled(s.Name); err != nil {
				return err
			} else {
				if alreadyInstalled {
					if v, err := SurveyForceInstall(s.Name); err != nil {
						return err
					} else {
						if !v {
							fmt.Println("install canceled")
							return nil
						}
					}
				}
			}

			cfg := c
			if shell != "" {
				cfg.Shell = []string{shell}
			}

			return qpm.Execute(cfg, s, qpm.Install, bufio.NewWriter(os.Stdout), bufio.NewWriter(os.Stderr))
		},
	}

	installCmd.PersistentFlags().StringVarP(&aquiferPath, "aquifer", "a", "", "Aquifer directory path")
	installCmd.PersistentFlags().StringVarP(&shell, "shell", "s", "", "Shell")
	rootCmd.AddCommand(installCmd)
}

func SurveyForceInstall(name string) (bool, error) {
	msg := name + " is already installed. Do you want to force installation?"

	v, err := survey.AskOneConfirm(msg, false)
	if err != nil {
		return false, err
	}

	return v, nil
}
