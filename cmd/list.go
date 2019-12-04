package cmd

import (
	"fmt"

	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

func newListCmd(cfg Config) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list all available applets in the repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			r := repo.New(cfg.RootDir)
			r.Init()

			for key, a := range r.Applets {
				fmt.Printf("%s:%s\n", key, a.Tag)
			}

			return nil
		},
	}
}
