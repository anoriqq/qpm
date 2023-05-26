package config

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/anoriqq/qpm"
	"github.com/goccy/go-yaml"
)

const (
	configFileName = "config"
	configFileType = "yml"
)

func InitConfigFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, ".config/qpm")

	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return "", err
	}

	configFilePath := filepath.Join(configPath, configFileName+"."+configFileType)
	if _, err := os.Stat(configFilePath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		if _, err := os.Create(configFilePath); err != nil {
			return "", err
		}
	}

	return configFilePath, nil
}

func ClearConfigFile() error {
	path, err := InitConfigFile()
	if err != nil {
		return err
	}

	if err := os.Truncate(path, 0); err != nil {
		return err
	}

	return nil
}

const (
	configNameAquiferPath    = "aquifer.path"
	configNameAquiferRemote  = "aquifer.remote"
	configNameGithubUsername = "github.username"
	configNameGithubToken    = "github.token"
	configNameShell          = "shell"
)

type rawConfig map[string]string

// WriteConfig 設定をLocalに保存する。
func WriteConfig(c qpm.Config, path string) error {
	r := make(rawConfig)
	r[configNameAquiferPath] = c.AquiferPath
	r[configNameAquiferRemote] = c.AquiferRemote.String()
	r[configNameGithubUsername] = c.GitHubUsername
	r[configNameGithubToken] = c.GitHubToken
	r[configNameShell] = c.Shell

	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := f.Truncate(0); err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		return err
	}

	return nil
}

// ReadConfig Localに保存してある設定を読み込む。
// 設定がなかった場合はerrorを返す。
func ReadConfig(path string) (qpm.Config, error) {
	c, err := readConfig(path)
	if err != nil {
		return qpm.Config{}, err
	}

	remoteURL, err := url.Parse(c[configNameAquiferRemote])
	if err != nil {
		return qpm.Config{}, err
	}

	return qpm.Config{
		AquiferPath:    c[configNameAquiferPath],
		AquiferRemote:  remoteURL,
		GitHubUsername: c[configNameGithubUsername],
		GitHubToken:    c[configNameGithubToken],
		Shell:          c[configNameShell],
	}, nil
}

func readConfig(path string) (rawConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return rawConfig{}, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return rawConfig{}, err
	}

	var c rawConfig
	if err := yaml.Unmarshal(b, &c); err != nil {
		return rawConfig{}, err
	}

	return c, nil
}

func StringConfig(path, match string) string {
	c, _ := readConfig(path)
	for key := range c {
		if !strings.HasPrefix(key, match) {
			delete(c, key)
		}
	}
	b, _ := yaml.Marshal(c)
	return string(b)
}
