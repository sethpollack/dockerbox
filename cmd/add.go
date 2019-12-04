package cmd

import (
	"github.com/sethpollack/dockerbox/registry"
	"github.com/sethpollack/dockerbox/repo"
	"github.com/spf13/cobra"
)

func newAddCmd(cfg Config) *cobra.Command {
	return &cobra.Command{
		Use:   "add",
		Short: "Add or update a repo in the registry.",
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.New(cfg.RootDir)
			if err != nil {
				return err
			}

			reg.Add(args[0], args[1])
			reg.Save()

			r := repo.New(cfg.RootDir)

			r.Update(reg)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
