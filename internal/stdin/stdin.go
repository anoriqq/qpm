package stdin

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"

	"github.com/anoriqq/qpm/internal/survey"
)

func SurveyAquiferRepoURL() (*url.URL, error) {
	msg := "Please enter qpm aquifer repository full URL."

	rawAquiferRepoURL, err := survey.AskOneInputRequired(msg, "")
	if err != nil {
		return nil, err
	}

	aquiferRepoURL, err := url.Parse(rawAquiferRepoURL)
	if err != nil {
		return nil, err
	}

	return aquiferRepoURL, nil
}

func SurveyDeleteAquiferDir() error {
	msg := "Do you want to replace the aquifer dir and continue?"

	canDelete, err := survey.AskOneConfirm(msg, false)
	if err != nil {
		return err
	}

	if !canDelete {
		return errors.New("to clone the aquifer repository, the aquifer dir must be deleted")
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

func SurveyAquiferDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	def := filepath.Join(home, ".qpm")

	msg := "Please enter qpm aquifer path."

	aquiferDir, err := survey.AskOneInputRequired(msg, def)
	if err != nil {
		return "", err
	}

	return aquiferDir, nil
}
