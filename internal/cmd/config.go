package cmd

import (
	"errors"

	"github.com/anoriqq/qpm/internal/service/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "update config",
	Example: `  # set AquiferDir to ~/.qpm
  qpm config AquiferDir ~/.qpm`,
	RunE: configRun,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configRun(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		config.PrintConfig()
		return nil
	}

	configField, configValue, err := getConfigInput(args)
	if err != nil {
		return err
	}

	err = config.SetConfig(configField, configValue)
	if err != nil {
		return err
	}

	config.PrintConfig()
	return nil
}

func getConfigInput(args []string) (field, value string, err error) {
	if len(args) < 2 {
		return "", "", errors.New("config field and value is required")
	}
	if len(args) > 2 {
		return "", "", errors.New("too many arguments")
	}

	return args[0], args[1], nil
}
