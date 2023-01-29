package cmd

import (
	"github.com/sethpollack/dockerbox/applet"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *dockerbox.Config, root *applet.Root) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use: "dockerbox",
	}

	cmd.AddCommand(
		newInstallCmd(cfg, root),
		newUninstallCmd(cfg, root),
		newDebugCmd(root),
		newVersionCmd(),
	)

	return cmd, nil
}
