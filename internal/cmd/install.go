package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/anoriqq/qpm/internal/config"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install packages",
	RunE:  installRun,
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func installRun(_ *cobra.Command, args []string) error {
	if !config.HasScriptDir() {
		scriptDir, err := surveyScriptDir()
		if err != nil {
			return err
		}

		err = config.SetScriptDir(scriptDir)
		if err != nil {
			return err
		}
	}

	pkgName, err := getPkgName(args)
	if err != nil {
		return err
	}

	installScriptPath, err := getInstallScriptPaht(config.Cfg.ScriptDir, pkgName)
	if err != nil {
		return err
	}

	return execInstallScript(installScriptPath)
}

func surveyScriptDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var scriptDir string
	p := &survey.Input{
		Message: "Please enter qpm script path.",
		Default: filepath.Join(home, ".qpm"),
	}
	err = survey.AskOne(p, &scriptDir, survey.WithValidator(survey.Required))
	if err != nil {
		return "", err
	}

	return scriptDir, nil
}

func getPkgName(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("package name required")
	}

	return args[0], nil
}

func getInstallScriptPaht(scriptDir, pkgName string) (string, error) {
	installScriptPath, err := filepath.Abs(fmt.Sprintf("%s/%s/install.sh", scriptDir, pkgName))
	if err != nil {
		return "", err
	}

	return installScriptPath, nil
}

func execInstallScript(installScriptPath string) error {
	_, err := os.Stat(installScriptPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("install script not found: %s", installScriptPath)
	}

	o, err := exec.Command("/bin/sh", installScriptPath).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(o))

	return nil
}
