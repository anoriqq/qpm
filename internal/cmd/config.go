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

var jsonToStringReplacer = strings.NewReplacer("{", "", "}", "")

func configRun(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		printConfig()
		return nil
	}

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
	case "scriptrepourl":
		config.SetScriptRepoURL(configValue)
	case "githubusername":
		config.SetGitHubUsername(configValue)
	case "githubaccesstoken":
		config.SetGitHubAccessToken(configValue)
	default:
		return fmt.Errorf("unknown config field: %s", configField)
	}

	printConfig()

	return nil
}

func printConfig() {
	cfgText := fmt.Sprintf("%+v\n", config.Cfg)
	cfgTexts := strings.Split(cfgText, " ")
	for _, t := range cfgTexts {
		text := strings.Replace(jsonToStringReplacer.Replace(t), ":", ": ", 1)
		fmt.Println(text)
	}
}
