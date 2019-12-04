package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/io"
	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

func newInstallCmd(cfg Config) *cobra.Command {
	var installAll bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install docker applet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !installAll && cfg.AppletName == "" {
				return fmt.Errorf("Missing --applet flag")
			}

			if err := io.EnsureDir(cfg.InstallDir); err != nil {
				return err
			}

			if installAll {
				r := repo.New(cfg.RootDir)
				r.Init()

				for key, _ := range r.Applets {
					install(cfg.DockerboxExe, cfg.InstallDir, key)
				}
			} else {
				install(cfg.DockerboxExe, cfg.InstallDir, cfg.AppletName)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&cfg.AppletName, "applet", "i", "", "applet to install")
	cmd.Flags().BoolVarP(&installAll, "all", "a", false, "install all applets.")

	return cmd
}

func install(exe, installDir, image string) {
	args := []string{
		"-s",
		"-f",
		exe,
		fmt.Sprintf("%s/%s", installDir, image),
	}

	cmd := exec.Command("ln", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
