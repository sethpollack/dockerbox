package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/applet"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/spf13/cobra"
)

func newUninstallCmd(cfg *dockerbox.Config, root *applet.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall docker applet",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, applet := range root.Applets {
				if _, ok := root.Ignore[applet.AppletName]; !ok {
					err := uninstall(cfg.InstallDir, applet.AppletName)
					if err != nil {
						return fmt.Errorf("failed to install %s: %v", applet.AppletName, err)
					}
				}
			}

			return nil
		},
	}

	return cmd
}

func uninstall(installDir, image string) error {
	args := []string{
		fmt.Sprintf("%s/%s", installDir, image),
	}
	cmd := exec.Command("rm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
