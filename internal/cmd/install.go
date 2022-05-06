package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

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
	pkgName, err := getPkgName(args)
	if err != nil {
		return err
	}

	installScriptPath, err := getInstallScriptPaht(pkgName)
	if err != nil {
		return err
	}

	return execInstallScript(installScriptPath)
}

func getPkgName(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("package name required")
	}

	return args[0], nil
}

func getInstallScriptPaht(pkgName string) (string, error) {
	// TODO: パスを決め打ちしないでconfigで持つ
	installScriptPath, err := filepath.Abs(fmt.Sprintf("./%s/install.sh", pkgName))
	if err != nil {
		return "", err
	}

	return installScriptPath, nil
}

func execInstallScript(installScriptPath string) error {
	o, err := exec.Command("/bin/sh", installScriptPath).Output()
	if err != nil {
		return err
	}

	fmt.Println(string(o))

	return nil
}
