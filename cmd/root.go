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

func Execute(conf Config) {
	cfg = conf
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
