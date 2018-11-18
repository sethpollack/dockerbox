package cmd

import (
	"fmt"

	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available applets in the repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := repo.New(cfg.RootDir)
		r.Init()

		for key, _ := range r.Applets {
			fmt.Println(key)
		}

		return nil
	},
}
