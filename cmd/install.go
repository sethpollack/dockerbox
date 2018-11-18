package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/sethpollack/dockerbox/io"
	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

var installAll bool

func init() {
	installCmd.Flags().StringVarP(&cfg.AppletName, "applet", "i", "", "applet to install")
	installCmd.Flags().BoolVarP(&installAll, "all", "a", false, "install all applets.")
}

var installCmd = &cobra.Command{
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
				install(key)
			}
		} else {
			install(cfg.AppletName)
		}

		return nil
	},
}

func install(image string) {
	cmd := linkCmd(image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func linkCmd(image string) *exec.Cmd {
	args := []string{
		"-s",
		"-f",
		cfg.DockerboxExe,
		fmt.Sprintf("%s/%s", cfg.InstallDir, image),
	}

	return exec.Command("ln", args...)
}
