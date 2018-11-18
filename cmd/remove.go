package cmd

import (
	"github.com/sethpollack/dockerbox/registry"
	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a repo from the registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := registry.New(cfg.RootDir)
		if err != nil {
			return err
		}

		reg.Remove(args[0])
		reg.Save()

		r := repo.New(cfg.RootDir)

		r.Update(reg)
		if err != nil {
			return err
		}

		return nil
	},
}
