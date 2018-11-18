package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	registryCmd.AddCommand(addCmd)
	registryCmd.AddCommand(removeCmd)
}

var registryCmd = &cobra.Command{
	Use: "registry",
}
