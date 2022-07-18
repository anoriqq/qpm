package aquifer

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

const (
	envOS   = "QPM_OS"
	envArch = "QPM_ARCH"
)

var headerOutput = color.New(color.FgHiCyan).Add(color.Bold)

type Plan struct {
	Dependencies []string
	Run          []string
}

type Aquifer struct {
	Version   string
	Name      string
	Install   map[string]Plan
	Uninstall map[string]Plan
}

func Install(aquiferDir, pkgName string) error {
	err := config.SetCfgWithSurvey(
		config.SetAquiferDirWithSurvey,
	)
	if err != nil {
		return err
	}

	aquiferPath, err := getAquiferPath(aquiferDir, pkgName)
	if err != nil {
		return err
	}

	aquifer, err := getAquifer(aquiferPath)
	if err != nil {
		return err
	}

	plan, err := getPlan(aquifer, runtime.GOOS)
	if err != nil {
		return err
	}

	return installAquifer(plan)
}

func getAquiferPath(aquiferDir, pkgName string) (string, error) {
	aquiferPath, err := filepath.Abs(fmt.Sprintf("%s/%s/latest.yml", aquiferDir, pkgName))
	if err != nil {
		return "", err
	}

	return aquiferPath, nil
}

func getAquifer(aquiferPath string) (Aquifer, error) {
	_, err := os.Stat(aquiferPath)
	if os.IsNotExist(err) {
		return Aquifer{}, fmt.Errorf("aquifer not found: %s", aquiferPath)
	}

	f, err := os.Open(aquiferPath)
	if err != nil {
		return Aquifer{}, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return Aquifer{}, err
	}

	var aquifer Aquifer
	err = yaml.Unmarshal(bytes, &aquifer)
	if err != nil {
		return Aquifer{}, err
	}

	return aquifer, nil
}

func getPlan(aquifer Aquifer, os string) (Plan, error) {
	plan, ok := aquifer.Install[os]
	if !ok {
		return Plan{}, fmt.Errorf("not declared os: %s", runtime.GOOS)
	}

	return plan, nil
}

func installAquifer(plan Plan) error {
	for i, v := range plan.Run {
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
