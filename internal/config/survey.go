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

func SetScriptDirWithSurvey() error {
	if hasScriptDir() {
		return nil
	}

	scriptDir, err := stdin.SurveyScriptDir()
	if err != nil {
		return err
	}

	err = SetScriptDir(scriptDir)
	if err != nil {
		return err
	}

	return nil
}

func SetScriptRepoURLWithSurvey() error {
	if hasScriptRepoURL() {
		return nil
	}

	scriptRepoURL, err := stdin.SurveyScriptRepoURL()
	if err != nil {
		return err
	}

	err = SetScriptRepoURL(scriptRepoURL.String())
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
