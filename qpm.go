package qpm

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func Version() string {
	return "v2.0.0"
}

type Config struct {
	AquiferPath   string
	AquiferRemote *url.URL

	GitHubUsername string
	GitHubToken    string

	Shell string
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
	step struct {
		name string
		run  string
	}
	job struct {
		dependencies []string
		steps        []step
	}
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
		OS           []string
		Dependencies []string
		Steps        []any
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

				slices.Sort(j.Dependencies)
				if len(slices.Compact(j.Dependencies)) != len(j.Dependencies) {
					return stratum{}, errors.New("deplicate packages in dependencies")
				}

				steps := make([]step, 0)
				for i, s := range j.Steps {
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

				s.Plan[action][os] = job{
					dependencies: j.Dependencies,
					steps:        steps,
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
func Execute(c Config, s stratum, action Action) error {
	cmd := exec.Command(c.Shell)

	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("QPM_OS=%s", runtime.GOOS),
		fmt.Sprintf("QPM_ARCH=%s", runtime.GOARCH),
	)

	cmd.Stdout = bufio.NewWriter(os.Stdout)
	cmd.Stderr = bufio.NewWriter(os.Stderr)

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

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	for i, r := range j.steps {
		io.WriteString(stdin, fmt.Sprintf(`echo "\e[96m[%d/%d] %s\e[m"`+"\n", i+1, len(j.steps), shellEscapeReplacer.Replace(r.name)))
		io.WriteString(stdin, fmt.Sprintln(r.run))
		io.WriteString(stdin, fmt.Sprintln(`if [ "$?" != 0 ]; then exit 1; fi`))
	}
	io.WriteString(stdin, fmt.Sprintln("echo '\\e[96m[Complete]\\e[m'"))

	if err := cmd.Start(); err != nil {
		return err
	}

	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
