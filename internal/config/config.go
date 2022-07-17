package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yml"
)

func InitConfig() error {
	configHome, err := getConfigHomeDir()
	if err != nil {
		return err
	}

	err = createIfNotExistConfigFile(configHome, configName, configType)
	if err != nil {
		return err
	}

	viper.AddConfigPath(configHome)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	return Load()
}

func getConfigHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config/qpm"), nil
}

func createIfNotExistConfigFile(configHome, configName, configType string) error {
	err := os.MkdirAll(configHome, os.ModePerm)
	if err != nil {
		return err
	}

	configPath := filepath.Join(configHome, configName+"."+configType)

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		_, err := os.Create(configPath)
		if err != nil {
			return err
		}
	}

	return nil
}

var Cfg Config

type Config struct {
	ScriptDir         string
	ScriptRepoURL     string
	GitHubUsername    string
	GitHubAccessToken string
}

func Load() error {
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		return err
	}

	return nil
}

func hasScriptDir() bool {
	return len(Cfg.ScriptDir) != 0
}

func SetScriptDir(scriptDir string) error {
	viper.Set("ScriptDir", scriptDir)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Cfg)
}

func hasScriptRepoURL() bool {
	return len(Cfg.ScriptRepoURL) != 0
}

func SetScriptRepoURL(scriptRepoURL string) error {
	viper.Set("ScriptRepoURL", scriptRepoURL)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Cfg)
}

func hasGitHubUsername() bool {
	return len(Cfg.GitHubUsername) != 0
}

func SetGitHubUsername(githubUsername string) error {
	viper.Set("GitHubUsername", githubUsername)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Cfg)
}

func hasGitHubAccessToken() bool {
	return len(Cfg.GitHubAccessToken) != 0
}

func SetGitHubAccessToken(githubAccessToken string) error {
	viper.Set("GitHubAccessToken", githubAccessToken)

	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Cfg)
}
