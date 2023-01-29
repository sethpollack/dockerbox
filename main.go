package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sethpollack/dockerbox/cmd"
	"github.com/sethpollack/dockerbox/cue"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/sethpollack/dockerbox/runner"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	// used to symlink to dockerbox binaries
	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("failed to get executable: %v", err)
		os.Exit(1)
	}
	// used to find local config overrides
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to get working directory: %v", err)
		os.Exit(1)
	}

	cfg, err := dockerbox.New(
		filepath.Base(os.Args[0]),
		wd,
		exe,
		os.Args[1:],
	)
	if err != nil {
		fmt.Printf("failed to create config: %v", err)
		os.Exit(1)
	}

	files, err := dockerbox.GetConfigurations(fs, wd, cfg.RootDir)
	if err != nil {
		fmt.Printf("failed to get configurations: %v", err)
		os.Exit(1)
	}

	root, err := cue.New(fs, files)
	if err != nil {
		fmt.Printf("failed to load applets: %v", err)
		os.Exit(1)
	}

	if _, err := fs.Stat(cfg.InstallDir); os.IsNotExist(err) {
		err := fs.MkdirAll(cfg.InstallDir, os.FileMode(0744))
		if err != nil {
			fmt.Printf("failed to create install directory: %v", err)
			os.Exit(1)
		}
	}

	switch cfg.EntryPoint {
	case "dockerbox":
		command, err := cmd.NewRootCmd(cfg, root)
		if err != nil {
			fmt.Printf("failed to create root command: %v", err)
			os.Exit(1)
		}

		err = command.Execute()
		if err != nil {
			fmt.Printf("failed to run command: %v", err)
			os.Exit(1)
		}
	default:
		cmds, err := root.Compile(cfg)
		if err != nil {
			fmt.Printf("failed to compile applet: %v", err)
			os.Exit(1)
		}

		err = runner.RunCmds(cmds)
		if err != nil {
			exiterr, ok := err.(*exec.ExitError)
			if ok {
				os.Exit(exiterr.ExitCode())
			}
			os.Exit(1)
		}
	}
}
