package aquifer

import (
	"fmt"
	"os"
	"time"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/anoriqq/qpm/internal/git"
	"github.com/anoriqq/qpm/internal/stdin"
)

var fileFmt = time.Now().Format("20060102150405")

func Pull() error {
	err := config.SetCfgWithSurvey(
		config.SetAquiferDirWithSurvey,
		config.SetAquiferRepoURLWithSurvey,
		config.SetGitHubUsernameWithSurvey,
		config.SetGitHubAccessTokenWithSurvey,
	)
	if err != nil {
		return err
	}

	err = stdin.SurveyDeleteAquiferDir()
	if err != nil {
		return err
	}

	oldName := fmt.Sprintf("%s.old_%s", config.Cfg.AquiferDir, fileFmt)
	err = os.Rename(config.Cfg.AquiferDir, oldName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	c, err := git.NewClient(config.Cfg.GitHubUsername, config.Cfg.GitHubAccessToken)
	if err != nil {
		return err
	}

	fmt.Println(config.Cfg.AquiferRepoURL)

	err = c.Clone(config.Cfg.AquiferDir, config.Cfg.AquiferRepoURL)
	if err != nil {
		return err
	}

	return nil
}
