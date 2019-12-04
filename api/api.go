package api

import (
	"fmt"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/sethpollack/dockerbox/cmd"
	"github.com/sethpollack/dockerbox/repo"
)

type Api struct {
	rootDir    string
	installDir string
	separator  string
	exe        string
	entrypoint string
	args       []string
}

func New(rootDir, installDir, separator, exe, entrypoint string, args []string) *Api {
	return &Api{
		rootDir:    rootDir,
		installDir: installDir,
		separator:  separator,
		exe:        exe,
		entrypoint: entrypoint,
		args:       args,
	}
}

func (a *Api) Run() error {
	cfg := cmd.Config{
		RootDir:      a.rootDir,
		InstallDir:   a.installDir,
		DockerboxExe: a.exe,
	}

	r := repo.New(a.rootDir)
	err := r.Init()
	if err != nil {
		return fmt.Errorf("Error loading repo: %v", err)
	}

	switch a.entrypoint {
	case "dockerbox":
		return cmd.NewRootCmd(cfg).Execute()
	default:
		app, ok := r.Applets[a.entrypoint]
		if !ok {
			fmt.Printf("Command %s not found", a.entrypoint)
			return err
		}

		fs, err := gpflag.Parse(&app)
		if err != nil {
			return err
		}

		dArgs, aArgs := splitArgs(a.separator, a.args)
		err = fs.Parse(dArgs)
		if err != nil {
			return err
		}

		app.PreExec()

		err = exec(r, app, aArgs...)
		if err != nil {
			return err
		}
	}

	return nil
}

func exec(r *repo.Repo, a repo.Applet, args ...string) error {
	for _, dep := range a.BeforeHooks {
		d, ok := r.Applets[dep]
		if !ok {
			return fmt.Errorf("dependency %s not found", dep)
		}
		err := exec(r, d)
		if err != nil {
			return err
		}
	}

	a.Exec(args...)

	for _, dep := range a.AfterHooks {
		d, ok := r.Applets[dep]
		if !ok {
			return fmt.Errorf("dependency %s not found", dep)
		}
		err := exec(r, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func splitArgs(separator string, args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == separator {
			return args[:i], args[i+1:]
		}
	}
	return []string{}, args
}
