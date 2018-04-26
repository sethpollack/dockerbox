package repo

import (
	"fmt"
	"os"
	"os/exec"
)

const dockerExe = "/usr/local/bin/docker"

type Applet struct {
	Name       string   `yaml:"name"`
	Image      string   `yaml:"image"`
	Privileged bool     `yaml:"privileged"`
	Daemonize  bool     `yaml:"daemonize"`
	Restart    string   `yaml:"restart"`
	Tag        string   `yaml:"image_tag"`
	WorkDir    string   `yaml:"work_dir"`
	Volumes    []string `yaml:"volumes"`
	Ports      []string `yaml:"ports`
	Entrypoint string   `yaml:"entrypoint"`
	Command    []string `yaml:"command"`
	Env        []string `yaml:"environment"`
}

func (a *Applet) Exec(extra []string) error {
	cmd := a.RunCmd(extra)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *Applet) RunCmd(extra []string) *exec.Cmd {
	args := []string{
		"run",
		"--rm",
		"-i",
	}
	for _, e := range a.Env {
		args = append(args, "-e", os.ExpandEnv(e))
	}
	for _, v := range a.Volumes {
		args = append(args, "-v", os.ExpandEnv(v))
	}
	for _, p := range a.Ports {
		args = append(args, "-p", p)
	}
	if a.Privileged {
		args = append(args, "--privileged")
	}
	if a.Daemonize {
		args = append(args, "-d")
	}
	if a.Entrypoint != "" {
		args = append(args, "--entrypoint", a.Entrypoint)
	}
	if a.Restart != "" {
		args = append(args, "--restart", a.Restart)
	}

	args = append(args, "-w", a.WorkDir)
	args = append(args, "--name", a.Name)
	args = append(args, fmt.Sprintf("%s:%s", a.Image, a.Tag))
	if len(a.Command) != 0 && len(extra) == 0 {
		args = append(args, a.Command...)
	} else {
		args = append(args, extra...)
	}

	return exec.Command(dockerExe, args...)
}
