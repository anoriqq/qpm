package aquifer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
)

const (
	envOS   = "QPM_OS"
	envArch = "QPM_ARCH"
	envEnv  = "QPM_ENV"
)

var (
	headerOutput = color.New(color.FgHiCyan).Add(color.Bold)
	errorOutput  = color.New(color.FgHiRed)
)

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
		return Plan{}, fmt.Errorf("not declared os: %s", os)
	}

	return plan, nil
}

type envMap map[string]string

func (e envMap) ToEnvSlice() []string {
	result := make([]string, len(e))

	var i int
	for k, v := range e {
		result[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	return result
}
func (e envMap) Append(key, val string) {
	e[key] = val
}

func installAquifer(plan Plan) error {
	envFilePath := path.Join(config.Cfg.AquiferDir, "tmp.env")

	f, err := os.Create(envFilePath)
	if err != nil {
		return err
	}

	defer f.Close()
	defer os.Remove(envFilePath)

	envs := envMap{
		envArch: runtime.GOARCH,
		envOS:   runtime.GOOS,
		envEnv:  envFilePath,
	}

	for i, v := range plan.Run {
		headerOutput.Printf("[%d/%d] %s\n", i+1, len(plan.Run), v)

		c := exec.Command("bash", "-c", v)

		c.Dir = config.Cfg.AquiferDir

		envFileBytes, err := io.ReadAll(f)
		if err != nil {
			return err
		}

		envStrings := strings.Split(string(envFileBytes), "\n")
		for _, v := range envStrings {
			keyval := strings.Split(v, "=")
			if len(keyval) == 2 {
				envs.Append(keyval[0], keyval[1])
			}
		}

		c.Env = append(c.Env, envs.ToEnvSlice()...)

		var stdout, stderr bytes.Buffer
		c.Stdout, c.Stderr = &stdout, &stderr

		err = c.Run()
		if err != nil {
			fmt.Printf("%v", stdout.String())
			errorOutput.Printf("[error] %v", stderr.String())
			return err
		}

		fmt.Printf("%v", stdout.String())
	}

	headerOutput.Println("[complete]")

	return nil
}
