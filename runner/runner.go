package runner

import (
	"os"
	"os/exec"
)

type Cmd struct {
	Silent bool
	Args   []string
}

func RunCmds(cmds []Cmd) error {
	for _, cmd := range cmds {
		exec := exec.Command(cmd.Args[0], cmd.Args[1:]...)
		if !cmd.Silent {
			exec.Stdout = os.Stdout
			exec.Stderr = os.Stderr
			exec.Stdin = os.Stdin
		}

		err := exec.Run()
		if err != nil && !cmd.Silent {
			return err
		}
	}

	return nil
}
