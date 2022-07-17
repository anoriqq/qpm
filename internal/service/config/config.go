package config

import (
	"fmt"
	"strings"

	"github.com/anoriqq/qpm/internal/config"
)

var jsonToStringReplacer = strings.NewReplacer("{", "", "}", "")

func PrintConfig() {
	cfgText := fmt.Sprintf("%+v\n", config.Cfg)
	cfgTexts := strings.Split(cfgText, " ")
	for _, t := range cfgTexts {
		text := strings.Replace(jsonToStringReplacer.Replace(t), ":", ": ", 1)
		fmt.Println(text)
	}
}

func SetConfig(configField, configValue string) error {
	switch strings.ToLower(configField) {
	case "scriptdir":
		config.SetScriptDir(configValue)
		return nil
	case "scriptrepourl":
		config.SetScriptRepoURL(configValue)
		return nil
	case "githubusername":
		config.SetGitHubUsername(configValue)
		return nil
	case "githubaccesstoken":
		config.SetGitHubAccessToken(configValue)
		return nil
	default:
		return fmt.Errorf("unknown config field: %s", configField)
	}
}
