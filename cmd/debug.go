package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/sethpollack/dockerbox/applet"
	"github.com/spf13/cobra"
)

func newDebugCmd(root *applet.Root) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "debug config files",
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := json.MarshalIndent(root, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal root: %v", err)
			}

			fmt.Println(string(bytes))

			return nil
		},
	}

	return cmd
}
