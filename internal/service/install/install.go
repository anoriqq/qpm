package install

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/anoriqq/qpm/internal/config"
)

func Install(pkgName string) error {
	err := config.SetCfgWithSurvey(
		config.SetScriptDirWithSurvey,
	)
	if err != nil {
		return err
	}

	installScriptPath, err := getInstallScriptPaht(config.Cfg.ScriptDir, pkgName)
	if err != nil {
		return err
	}

	return execInstallScript(installScriptPath)
}

func getInstallScriptPaht(scriptDir, pkgName string) (string, error) {
	installScriptPath, err := filepath.Abs(fmt.Sprintf("%s/%s/latest.sh", scriptDir, pkgName))
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

	c := exec.Command("/bin/sh", installScriptPath, "install", runtime.GOOS, runtime.GOARCH)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	c.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return c.Wait()
}
