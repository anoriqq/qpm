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
	ScriptDir string
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

func HasScriptDir() bool {
	return len(Cfg.ScriptDir) != 0
}

func SetScriptDir(scriptDir string) error {
	viper.Set("ScriptDir", scriptDir)

	err :=  viper.WriteConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&Cfg)
}
