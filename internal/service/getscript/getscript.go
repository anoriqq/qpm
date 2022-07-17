package getscript

import (
	"fmt"
	"os"
	"time"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/anoriqq/qpm/internal/git"
	"github.com/anoriqq/qpm/internal/stdin"
)

func GetScript() error {
	err := config.SetCfgWithSurvey(
		config.SetScriptDirWithSurvey,
		config.SetScriptRepoURLWithSurvey,
		config.SetGitHubUsernameWithSurvey,
		config.SetGitHubAccessTokenWithSurvey,
	)
	if err != nil {
		return err
	}

	err = stdin.SurveyDeleteScriptDir()
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
