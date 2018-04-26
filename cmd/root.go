package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	appletName      string
	dockerboxExe, _ = os.Executable()
	prefix          = os.ExpandEnv("$HOME/.dockerbox/bin")
	rootCmd         = &cobra.Command{
		Use: "dockerbox",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
