package cmd

import (
	"fmt"

	"github.com/sethpollack/dockerbox/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version %s\nCommit %s\n", version.Version, version.Commit)
		},
	}
}
