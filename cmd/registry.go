package cmd

import (
	"github.com/spf13/cobra"
)

func newRegistryCmd(cfg Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "registry",
	}
	cmd.AddCommand(
		newAddCmd(cfg),
		newRemoveCmd(cfg),
	)
	return cmd
}
