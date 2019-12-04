package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

func newUninstallCmd(cfg Config) *cobra.Command {
	var uninstallAll bool
	cmd := &cobra.Command{
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
					uninstall(cfg.InstallDir, key)
				}
			} else {
				uninstall(cfg.InstallDir, cfg.AppletName)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&cfg.AppletName, "applet", "i", "", "applet to uninstall")
	cmd.Flags().BoolVarP(&uninstallAll, "all", "a", false, "uninstall all applets.")

	return cmd
}

func uninstall(installDir, image string) {
	args := []string{
		fmt.Sprintf("%s/%s", installDir, image),
	}
	cmd := exec.Command("rm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
