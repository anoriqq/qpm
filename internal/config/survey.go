package config

import (
	"github.com/anoriqq/qpm/internal/stdin"
)

func SetCfgWithSurvey(setters ...func() error) error {
	for _, setter := range setters {
		err := setter()
		if err != nil {
			return err
		}
	}

	return nil
}

func SetAquiferDirWithSurvey() error {
	if hasAquiferDir() {
		return nil
	}

	aquiferDir, err := stdin.SurveyAquiferDir()
	if err != nil {
		return err
	}

	err = SetAquiferDir(aquiferDir)
	if err != nil {
		return err
	}

	return nil
}

func SetAquiferRepoURLWithSurvey() error {
	if hasAquiferRepoURL() {
		return nil
	}

	aquiferRepoURL, err := stdin.SurveyAquiferRepoURL()
	if err != nil {
		return err
	}

	err = SetAquiferRepoURL(aquiferRepoURL.String())
	if err != nil {
		return err
	}

	return nil
}

func SetGitHubUsernameWithSurvey() error {
	if hasGitHubUsername() {
		return nil
	}

	githubUsername, err := stdin.SurveyGitHubUsername()
	if err != nil {
		return err
	}

	err = SetGitHubUsername(githubUsername)
	if err != nil {
		return err
	}

	return nil
}

func SetGitHubAccessTokenWithSurvey() error {
	if hasGitHubAccessToken() {
		return nil
	}

	githubAccessToken, err := stdin.SurveyGitHubAccessToken()
	if err != nil {
		return err
	}

	err = SetGitHubAccessToken(githubAccessToken)
	if err != nil {
		return err
	}

	return nil
}
