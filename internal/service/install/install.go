package install

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
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

	return installAquifer(installScriptPath)
}

func getInstallScriptPaht(scriptDir, pkgName string) (string, error) {
	installScriptPath, err := filepath.Abs(fmt.Sprintf("%s/%s/latest.yml", scriptDir, pkgName))
	if err != nil {
		return "", err
	}

	return installScriptPath, nil
}

type plan struct {
	Dependencies []string
	Run          []string
}

type Aquifer struct {
	Version   string
	Name      string
	Install   map[string]plan
	Uninstall map[string]plan
}

var headerOutput = color.New(color.FgHiCyan).Add(color.Bold)

const (
	envOS   = "QPM_OS"
	envArch = "QPM_ARCH"
)

func installAquifer(installScriptPath string) error {
	_, err := os.Stat(installScriptPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("install script not found: %s", installScriptPath)
	}

	f, err := os.Open(installScriptPath)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var aquifer Aquifer
	err = yaml.Unmarshal(bytes, &aquifer)
	if err != nil {
		return err
	}

	cmds, ok := aquifer.Install[runtime.GOOS]
	if !ok {
		return fmt.Errorf("not declared os: %s", runtime.GOOS)
	}

	for i, v := range cmds.Run {
		headerOutput.Printf("[%d] %s\n", i+1, v)

		c := exec.Command("bash", "-c", v)
		c.Env = append(c.Env,
			fmt.Sprintf("%s=%s", envOS, runtime.GOARCH),
			fmt.Sprintf("%s=%s", envArch, runtime.GOOS),
		)

		output, err := c.Output()
		if err != nil {
			return err
		}

		fmt.Printf("%v", string(output))
	}

	return nil
}
