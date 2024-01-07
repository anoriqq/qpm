package qpm

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func Version() string {
	return "v0.0.14"
}

type Config struct {
	AquiferPath   string
	AquiferRemote *url.URL

	GitHubUsername string
	GitHubToken    string

	Shell []string
}

// PullAquifer AquiferをRemoteから取得する。
func PullAquifer(ctx context.Context, c Config) error {
	return nil
}

type (
	Action string
	OS     string
)

const (
	Install   Action = "install"
	Uninstall Action = "uninstall"

	linux  OS = "linux"
	darwin OS = "darwin"
)

var Actions = map[string]Action{
	string(Install):   Install,
	string(Uninstall): Uninstall,
}

var OSs = map[string]OS{
	string(linux):  linux,
	string(darwin): darwin,
}

func parseAction(v string) (Action, error) {
	if a, ok := Actions[v]; ok {
		return a, nil
	}

	return "", errors.Errorf("unknown Action v=%q", v)
}

func parseOS(v string) (OS, error) {
	if os, ok := OSs[v]; ok {
		return os, nil
	}

	return "", errors.Errorf("unknown OS v=%q", v)
}

type (
	osToJob map[OS]job
	plan    map[Action]osToJob
	stratum struct {
		Plan plan
		Name string
	}
)

// ReadStratum AquiferPathにあるStratumのうち、指定されたStratumを取得する。
func ReadStratum(c Config, name string) (stratum, error) {
	path := os.ExpandEnv(c.AquiferPath)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return stratum{}, fmt.Errorf("aquifer not found in %s", path)
	}

	stratumPath, err := filepath.Abs(fmt.Sprintf("%s/%s.yml", path, name))
	if err != nil {
		return stratum{}, err
	}

	if _, err := os.Stat(stratumPath); os.IsNotExist(err) {
		return stratum{}, fmt.Errorf("stratum not found path=%s", stratumPath)
	}

	f, err := os.Open(stratumPath)
	if err != nil {
		return stratum{}, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return stratum{}, err
	}

	var ay map[string][]struct {
		OS         []string
		Shell      []string
		Dependency []string
		Step       []any
	}
	if err := yaml.Unmarshal(b, &ay); err != nil {
		fmt.Println(yaml.FormatError(err, true, true))
		return stratum{}, err
	}

	s := stratum{
		Plan: make(plan),
		Name: name,
	}
	for actionStr, jobs := range ay {
		action, err := parseAction(actionStr)
		if err != nil {
			return stratum{}, err
		}

		if s.Plan[action] == nil {
			s.Plan[action] = make(osToJob)
		}

		for _, j := range jobs {
			for _, osStr := range j.OS {
				os, err := parseOS(osStr)
				if err != nil {
					return stratum{}, err
				}

				slices.Sort(j.Dependency)
				if len(slices.Compact(j.Dependency)) != len(j.Dependency) {
					return stratum{}, errors.New("duplicate packages in dependencies")
				}

				steps := make([]step, 0)
				for i, s := range j.Step {
					switch v := s.(type) {
					case string:
						steps = append(steps, step{
							name: v,
							run:  v,
						})
					case map[string]any:
						n, ok := v["name"]
						if !ok {
							return stratum{}, errors.Errorf("invalid value action=%v os=%v step-index=%d", action, os, i)
						}
						r, ok := v["run"]
						if !ok {
							return stratum{}, errors.Errorf("invalid value action=%v os=%v step-index=%d", action, os, i)
						}

						nn, ok := n.(string)
						if !ok {
							return stratum{}, errors.Errorf("invalid value action=%v os=%v step-index=%d", action, os, i)
						}
						rr, ok := r.(string)
						if !ok {
							return stratum{}, errors.Errorf("invalid value action=%v os=%v step-index=%d", action, os, i)
						}

						steps = append(steps, step{
							name: nn,
							run:  rr,
						})
					default:
						return stratum{}, errors.Errorf("invalid value action=%v os=%v step-index=%d", action, os, i)
					}
				}

				shell := make(map[string]struct{}, len(j.Shell))
				for _, s := range j.Shell {
					shell[s] = struct{}{}
				}

				s.Plan[action][os] = job{
					dependency:     j.Dependency,
					availableShell: shell,
					step:           steps,
				}
			}
		}
	}

	return s, nil
}

var shellEscapeReplacer = strings.NewReplacer("$", `\$`)

var ErrPackageAlreadyInstalled = errors.New("package already installed")

func IsAlreadyInstalled(s stratum) (bool, error) {
	path, err := exec.LookPath(s.Name)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		} else {
			return false, err
		}
	}
	if len(path) == 0 {
		return false, nil
	}

	return true, nil
}

// Execute aquiferを実行する。
func Execute(c Config, s stratum, action Action, stdout, stderr io.Writer) error {
	a, ok := s.Plan[action]
	if !ok {
		return errors.Errorf("%s is an Action not defined in the stratum", action)
	}

	os, err := parseOS(runtime.GOOS)
	if err != nil {
		return err
	}

	j, ok := a[os]
	if !ok {
		return errors.Errorf("%s is an OS not defined in the stratum", action)
	}

	shell, ok := j.shell(c.Shell)
	if !ok {
		return errors.Errorf("%s is a shell not defined in the stratum", c.Shell)
	}

	cmd := exec.Command(shell)

	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("QPM_OS=%s", runtime.GOOS),
		fmt.Sprintf("QPM_ARCH=%s", runtime.GOARCH),
	)

	cmd.Stdout, cmd.Stderr = stdout, stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	header := fmt.Sprintf("=> %s %s by %s", strings.ToUpper(string(action[0]))+string(action)[1:], s.Name, shell)
	echo(stdin, color.FgYellow, header)
	for i, r := range j.step {
		title := fmt.Sprintf("[%d/%d] %s", i+1, len(j.step), shellEscapeReplacer.Replace(r.name))
		echo(stdin, color.FgCyan, title)
		fmt.Fprintln(stdin, r.run)
		fmt.Fprintln(stdin, `if [ "$?" != 0 ]; then exit 1; fi`)
	}
	echo(stdin, color.FgYellow, "=> Complete")

	if err := cmd.Start(); err != nil {
		return err
	}

	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func echo(w io.Writer, att color.Attribute, s string) {
	str := color.New(att).Sprint(s)
	fmt.Fprintf(w, "echo '%s'\n", str)
}
