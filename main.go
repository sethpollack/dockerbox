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
		fmt.Printf("Error loading repo %s", entrypoint)
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

		a.Exec(args)
	}
}
