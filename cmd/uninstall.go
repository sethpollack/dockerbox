package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

var uninstallAll bool

func init() {
	rootCmd.AddCommand(uninstallCmd)

	uninstallCmd.Flags().StringVarP(&cfg.AppletName, "applet", "i", "", "applet to uninstall")
	uninstallCmd.Flags().BoolVarP(&uninstallAll, "all", "a", false, "uninstall all applets.")
}

var (
	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall docker applet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !uninstallAll && cfg.AppletName == "" {
				return fmt.Errorf("Missing --applet flag")
			}

			if uninstallAll {
				r := repo.New(cfg.RootDir)
				r.Init()

				for key, _ := range r.Applets {
					uninstall(key)
				}
			} else {
				uninstall(cfg.AppletName)
			}

			return nil
		},
	}
)

func uninstall(image string) {
	cmd := unlinkCmd(image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func unlinkCmd(image string) *exec.Cmd {
	args := []string{
		fmt.Sprintf("%s/%s", cfg.InstallDir, image),
	}

	return exec.Command("rm", args...)
}
