package cmd

import (
	"github.com/spf13/cobra"
)

type Config struct {
	AppletName   string
	RootDir      string
	InstallDir   string
	DockerboxExe string
}

func NewRootCmd(cfg Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "dockerbox",
	}
	cmd.AddCommand(
		newUpdateCmd(cfg),
		newListCmd(cfg),
		newInstallCmd(cfg),
		newUninstallCmd(cfg),
		newRegistryCmd(cfg),
		newVersionCmd(cfg),
	)
	return cmd
}
