package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	AppletName   string
	RootDir      string
	InstallDir   string
	DockerboxExe string
}

var (
	cfg     Config
	rootCmd = &cobra.Command{
		Use: "dockerbox",
	}
)

func init() {
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(registryCmd)
}

func Execute(conf Config) {
	cfg = conf
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
