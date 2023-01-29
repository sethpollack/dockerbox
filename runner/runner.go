package runner

import (
	"os"
	"os/exec"
)

func RunCmds(cmds [][]string) error {
	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
