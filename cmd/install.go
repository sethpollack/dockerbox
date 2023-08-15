package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/applet"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/spf13/cobra"
)

func newInstallCmd(cfg *dockerbox.Config, root *applet.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install docker applet",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, applet := range root.Applets {
				if _, ok := root.Ignore[applet.AppletName]; !ok {
					err := install(cfg.DockerboxExe, cfg.InstallDir, applet.AppletName)
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

func install(exe, installDir, image string) error {
	args := []string{
		"-s",
		"-f",
		exe,
		fmt.Sprintf("%s/%s", installDir, image),
	}

	cmd := exec.Command("ln", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
