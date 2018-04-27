package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethpollack/dockerbox/cmd"
	"github.com/sethpollack/dockerbox/repo"
)

func main() {
	entrypoint := filepath.Clean(os.Args[0])
	args := os.Args[1:]

	r := repo.New()
	err := r.Init()
	if err != nil {
		fmt.Printf("Error loading repo: %v", err)
		os.Exit(1)
	}

	switch entrypoint {
	case "dockerbox":
		cmd.Execute()
	default:
		a, ok := r.Applets[entrypoint]
		if !ok {
			fmt.Printf("Command %s not found", entrypoint)
			os.Exit(1)
		}

		err := Exec(r, a, args...)
		if err != nil {
			fmt.Printf("Error running applet: %v", err)
			os.Exit(1)
		}
	}
}

func Exec(r *repo.Repo, a repo.Applet, args ...string) error {
	for _, dep := range a.Dependencies {
		d, ok := r.Applets[dep]
		if !ok {
			return fmt.Errorf("dependency %s not found", dep)
		}
		err := Exec(r, d)
		if err != nil {
			return err
		}
	}

	return a.Exec(args...)
}
