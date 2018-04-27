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
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&appletName, "applet", "i", "", "applet to install")
	installCmd.Flags().BoolVarP(&installAll, "all", "a", false, "install all applets.")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "install docker applet",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !installAll && appletName == "" {
			return fmt.Errorf("Missing --applet flag")
		}

		if err := io.EnsureDir(prefix); err != nil {
			return err
		}

		if installAll {
			r := repo.New()
			r.Init()

			for _, a := range r.Applets {
				install(a.Name)
			}
		} else {
			install(appletName)
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
		dockerboxExe,
		fmt.Sprintf("%s/%s", prefix, image),
	}

	return exec.Command("ln", args...)
}
