package cmd

import (
	"errors"

	"github.com/anoriqq/qpm/internal/config"
	"github.com/anoriqq/qpm/internal/service/aquifer"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install packages",
	RunE:  installRun,
}

var aquiferDir string

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.PersistentFlags().StringVarP(&aquiferDir, "aquifer-dir", "d", "", "aquifer dir")
}

func installRun(_ *cobra.Command, args []string) error {
	pkgName, err := getPkgName(args)
	if err != nil {
		return err
	}

	return aquifer.Install(getAquiferDir(), pkgName)
}

func getPkgName(args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("package name required")
	}

	return args[0], nil
}

func getAquiferDir() string {
	if len(aquiferDir) != 0 {
		return aquiferDir
	}

	return config.Cfg.AquiferDir
}
