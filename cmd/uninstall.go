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

	uninstallCmd.Flags().StringVarP(&appletName, "applet", "i", "", "applet to uninstall")
	uninstallCmd.Flags().BoolVarP(&uninstallAll, "all", "a", false, "uninstall all applets.")
}

var (
	uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall docker applet",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !uninstallAll && appletName == "" {
				return fmt.Errorf("Missing --applet flag")
			}

			if uninstallAll {
				r := repo.New()
				r.Init()

				for _, a := range r.Applets {
					uninstall(a.Name)
				}
			} else {
				uninstall(appletName)
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
		fmt.Sprintf("%s/%s", prefix, image),
	}

	return exec.Command("rm", args...)
}
