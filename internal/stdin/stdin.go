package stdin

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/anoriqq/qpm/internal/survey"
)

func SurveyScriptRepoURL() (*url.URL, error) {
	msg := "Please enter qpm script repository full URL."

	rawScriptRepoURL, err := survey.AskOneInputRequired(msg, "")
	if err != nil {
		return nil, err
	}

	scriptRepoURL, err := url.Parse(rawScriptRepoURL)
	if err != nil {
		return nil, err
	}

	return scriptRepoURL, nil
}

func SurveyDeleteScriptDir() error {
	msg := "Do you want to replace the script dir and continue?"

	canDelete, err := survey.AskOneConfirm(msg, false)
	if err != nil {
		return err
	}

	if !canDelete {
		return errors.New("to clone the script repository, the script dir must be deleted")
	}

	return nil
}

func SurveyGitHubUsername() (string, error) {
	msg := "Please enter github username."

	githubUsername, err := survey.AskOneInputRequired(msg, "")
	if err != nil {
		return "", err
	}

	return githubUsername, nil
}

func SurveyGitHubAccessToken() (string, error) {
	msg := "Please enter github access token."

	githubAccessToken, err := survey.AskOnePasswordRequired(msg)
	if err != nil {
		return "", err
	}

	return githubAccessToken, nil
}

func SurveyScriptDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	def := filepath.Join(home, ".qpm")

	msg := "Please enter qpm script path."

	scriptDir, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return "", err
	}

	return scriptDir, nil
}
