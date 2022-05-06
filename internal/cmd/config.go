package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "update config",
	Example: `  # set ScriptDir to ~/.qpm
  qpm config ScriptDir ~/.qpm`,
	RunE: configRun,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

var r = strings.NewReplacer("{", "", "}", "", ":", ": ")

func configRun(_ *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("config field and value is required")
	}
	if len(args) != 2 {
		return errors.New("too many arguments")
	}

	configField, configValue := args[0], args[1]

	switch strings.ToLower(configField) {
	case "scriptdir":
		config.SetScriptDir(configValue)
	default:
		return fmt.Errorf("unknown config field: %s", configField)
	}

	fmt.Println(r.Replace(fmt.Sprintf("%+v\n", config.Cfg)))

	return nil
}
