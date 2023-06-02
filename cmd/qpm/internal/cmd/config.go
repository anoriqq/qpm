package cmd

import (
	"fmt"
	"net/url"

	"github.com/anoriqq/qpm"
	"github.com/anoriqq/qpm/cmd/qpm/internal/config"
	"github.com/anoriqq/qpm/cmd/qpm/internal/survey"
	"github.com/spf13/cobra"
)

func init() {
	var isInit bool
	var isClear bool

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage qpm config",
		Example: `  # Set aquifer.path to ~/.qpm
  qpm config aquifer.path ~/.qpm`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if isClear {
				return config.ClearConfigFile()
			}

			path, err := config.InitConfigFile()
			if err != nil {
				return err
			}

			if isInit {
				current, err := config.ReadConfig(path)
				if err != nil {
					return err
				}
				var c qpm.Config
				if v, err := SurveyAquiferPath(&current.AquiferPath); err != nil {
					return err
				} else {
					c.AquiferPath = v
				}
				if v, err := SurveyAquiferRemote(current.AquiferRemote); err != nil {
					return err
				} else {
					c.AquiferRemote = v
				}
				if v, err := SurveyGitHubUsername(&current.GitHubUsername); err != nil {
					return err
				} else {
					c.GitHubUsername = v
				}
				if v, err := SurveyGitHubToken(&current.GitHubToken); err != nil {
					return err
				} else {
					c.GitHubToken = v
				}
				if v, err := SurveyShell(&current.Shell); err != nil {
					return err
				} else {
					c.Shell = v
				}
				return config.WriteConfig(c, path)
			} else {
				switch len(args) {
				case 0:
					fmt.Print(config.StringConfig(path, ""))
					return nil
				case 1:
					fmt.Print(config.StringConfig(path, args[0]))
					return nil
				default:
					return nil
				}
			}
		},
	}

	configCmd.Flags().BoolVarP(&isInit, "init", "i", false, "Interactive initialisation")
	configCmd.Flags().BoolVarP(&isClear, "clear", "", false, "Clear config file")
	rootCmd.AddCommand(configCmd)
}

func SurveyAquiferPath(current *string) (string, error) {
	msg := "Please enter qpm aquifer path."

	def := "$HOME/.qpm"
	if current != nil && *current != "" {
		def = *current
	}

	v, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return "", err
	}

	return v, nil
}

func SurveyAquiferRemote(current *url.URL) (*url.URL, error) {
	msg := "Please enter qpm aquifer repository full URL."

	def := ""
	if current != nil && current.String() != "" {
		def = current.String()
	}

	raw, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return nil, err
	}

	v, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func SurveyGitHubUsername(current *string) (string, error) {
	msg := "Please enter GitHub username."

	def := ""
	if current != nil && *current != "" {
		def = *current
	}

	v, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return "", err
	}

	return v, nil
}

func SurveyGitHubToken(current *string) (string, error) {
	msg := "Please enter GitHub access token. If nothing is entered, the current config will be taken over."

	v, err := survey.AskOnePassword(msg)
	if err != nil {
		return "", err
	}

	if current != nil && v == "" {
		return *current, nil
	}

	return v, nil
}

func SurveyShell(current *string) (string, error) {
	msg := "Please enter shell command you want to use to execute stratum."

	def := "bash"
	if current != nil && *current != "" {
		def = *current
	}

	v, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return "", err
	}

	return v, nil
}
