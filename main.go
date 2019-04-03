package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethpollack/dockerbox/cmd"
	"github.com/sethpollack/dockerbox/io"
	"github.com/sethpollack/dockerbox/repo"
)

func main() {
	entrypoint := filepath.Base(os.Args[0])
	args := os.Args[1:]

	rootDir := io.GetEnv("DOCKERBOX_ROOT_DIR", "$HOME/.dockerbox")
	installDir := io.GetEnv("DOCKERBOX_INSTALL_DIR", rootDir+"/bin")

	exe, _ := os.Executable()

	cfg := cmd.Config{
		RootDir:      rootDir,
		InstallDir:   installDir,
		DockerboxExe: exe,
	}

	r := repo.New(rootDir)
	err := r.Init()
	if err != nil {
		fmt.Printf("Error loading repo: %v", err)
		os.Exit(1)
	}

	switch entrypoint {
	case "dockerbox":
		cmd.Execute(cfg)
	default:
		a, ok := r.Applets[entrypoint]
		if !ok {
			fmt.Printf("Command %s not found", entrypoint)
			os.Exit(1)
		}

		if a.Kill {
			a.PreExec()
		}

		err := Exec(r, a, args...)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}
}

func Exec(r *repo.Repo, a repo.Applet, args ...string) error {
	for _, dep := range a.BeforeHooks {
		d, ok := r.Applets[dep]
		if !ok {
			return fmt.Errorf("dependency %s not found", dep)
		}
		err := Exec(r, d)
		if err != nil {
			return err
		}
	}

	err := a.Exec(args...)
	if err != nil {
		return err
	}

	for _, dep := range a.AfterHooks {
		d, ok := r.Applets[dep]
		if !ok {
			return fmt.Errorf("dependency %s not found", dep)
		}
		err := Exec(r, d)
		if err != nil {
			return err
		}
	}

	return nil
}
