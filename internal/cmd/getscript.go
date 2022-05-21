package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anoriqq/qpm/internal/config"
	"github.com/anoriqq/qpm/internal/git"
	"github.com/spf13/cobra"
)

var getscriptCmd = &cobra.Command{
	Use:   "getscript",
	Short: "get script dir form remote repository",
	RunE:  getscriptRun,
}

func init() {
	rootCmd.AddCommand(getscriptCmd)
}

func getscriptRun(_ *cobra.Command, _ []string) error {
	if !config.HasScriptDir() {
		scriptDir, err := surveyScriptDir()
		if err != nil {
			return err
		}

		err = config.SetScriptDir(scriptDir)
		if err != nil {
			return err
		}
	}

	if !config.HasScriptRepoURL() {
		scriptRepoURL, err := surveyScriptRepoURL()
		if err != nil {
			return err
		}

		err = config.SetScriptRepoURL(scriptRepoURL.String())
		if err != nil {
			return err
		}
	}

	if !config.HasGitHubUsername() {
		githubUsername, err := surveyGitHubUsername()
		if err != nil {
			return err
		}

		err = config.SetGitHubUsername(githubUsername)
		if err != nil {
			return err
		}
	}

	if !config.HasGitHubAccessToken() {
		githubAccessToken, err := surveyGitHubAccessToken()
		if err != nil {
			return err
		}

		err = config.SetGitHubAccessToken(githubAccessToken)
		if err != nil {
			return err
		}
	}

	err := surveyDeleteScriptDir()
	if err != nil {
		return err
	}

	err = os.Rename(config.Cfg.ScriptDir, fmt.Sprintf("%s.old_%s", config.Cfg.ScriptDir, time.Now().Format("20060102150405")))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	c, err := git.NewClient(config.Cfg.GitHubUsername, config.Cfg.GitHubAccessToken)
	if err != nil {
		return err
	}

	fmt.Println(config.Cfg.ScriptRepoURL)

	err = c.Clone(config.Cfg.ScriptDir, config.Cfg.ScriptRepoURL)
	if err != nil {
		return err
	}

	return nil
}

func surveyScriptRepoURL() (*url.URL, error) {
	var rawScriptRepoURL string
	p := &survey.Input{
		Message: "Please enter qpm script repository full URL.",
	}
	err := survey.AskOne(p, &rawScriptRepoURL, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	scriptRepoURL, err := url.Parse(rawScriptRepoURL)
	if err != nil {
		return nil, err
	}

	return scriptRepoURL, nil
}

func surveyDeleteScriptDir() error {
	var canDelete bool
	p := &survey.Confirm{
		Message: "Do you want to replace the script dir and continue?",
		Default: false,
	}
	err := survey.AskOne(p, &canDelete)
	if err != nil {
		return err
	}

	if !canDelete {
		return errors.New("To clone the script repository, the script dir must be deleted.")
	}

	return nil
}

func surveyGitHubUsername() (string, error) {
	var githubUsername string
	p := &survey.Input{
		Message: "Please enter github username.",
	}
	err := survey.AskOne(p, &githubUsername, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return githubUsername, nil
}

func surveyGitHubAccessToken() (string, error) {
	var githubAccessToken string
	p := &survey.Password{
		Message: "Please enter github access token.",
	}
	err := survey.AskOne(p, &githubAccessToken, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return githubAccessToken, nil
}
