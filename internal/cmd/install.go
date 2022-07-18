package cmd

import (
	"errors"

	"github.com/anoriqq/qpm/internal/service/aquifer"
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

	return aquifer.Install(pkgName)
}

func getPkgName(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("package name required")
	}

	return args[0], nil
}
