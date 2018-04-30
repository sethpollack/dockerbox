package repo

import (
	"fmt"
	"os"
	"os/exec"
)

const dockerExe = "/usr/local/bin/docker"

type Applet struct {
	Name       string `yaml:"name,omitempty"`
	WorkDir    string `yaml:"work_dir,omitempty"`
	Entrypoint string `yaml:"entrypoint,omitempty"`
	Restart    string `yaml:"restart,omitempty"`
	Network    string `yaml:"network,omitempty"`

	RM          bool `yaml:"rm,omitempty"`
	TTY         bool `yaml:"tty,omitempty"`
	Interactive bool `yaml:"interactive,omitempty"`
	Privileged  bool `yaml:"privileged,omitempty"`
	Detach      bool `yaml:"detach,omitempty"`
	Kill        bool `yaml:"kill,omitempty"`

	Env          []string `yaml:"environment,omitempty"`
	Volumes      []string `yaml:"volumes,omitempty"`
	Ports        []string `yaml:"ports,omitempty"`
	EnvFile      []string `yaml:"env_file,omitempty"`
	Dependencies []string `yaml:"dependencies,omitempty"`
	Links        []string `yaml:"links,omitempty"`

	Image string `yaml:"image,omitempty"`
	Tag   string `yaml:"image_tag,omitempty"`

	Command []string `yaml:"command,omitempty"`
}

func (a *Applet) Exec(extra ...string) error {
	cmd := a.RunCmd(extra)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error runnng applet %s: %v", a.Name, err)
	}
	return nil
}

func (a *Applet) PreExec() {
	cmd := a.KillCmd()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("error killing %s: %v", a.Name, err)
	}
}

func (a *Applet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawApplet Applet
	raw := rawApplet{
		RM:          true,
		Interactive: true,
		TTY:         true,
		Tag:         "latest",
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*a = Applet(raw)
	return nil
}

func (a *Applet) KillCmd() *exec.Cmd {
	args := []string{
		"kill",
		a.Name,
	}

	return exec.Command(dockerExe, args...)
}

func (a *Applet) RunCmd(extra []string) *exec.Cmd {
	args := []string{
		"run",
	}

	if a.Name != "" {
		args = append(args, "--name", a.Name)
	}
	if a.WorkDir != "" {
		args = append(args, "--workdir", a.WorkDir)
	}
	if a.Entrypoint != "" {
		args = append(args, "--entrypoint", a.Entrypoint)
	}
	if a.Restart != "" {
		args = append(args, "--restart", a.Restart)
	}
	if a.Network != "" {
		args = append(args, "--network", a.Network)
	}

	if a.RM {
		args = append(args, "--rm")
	}
	if a.Interactive {
		args = append(args, "--interactive")
	}
	if a.Privileged {
		args = append(args, "--privileged")
	}
	if a.Detach {
		args = append(args, "--detach")
	}
	if a.TTY {
		args = append(args, "--tty")
	}

	for _, f := range a.Env {
		args = append(args, "-e", os.ExpandEnv(f))
	}
	for _, f := range a.Volumes {
		args = append(args, "-v", os.ExpandEnv(f))
	}
	for _, f := range a.Ports {
		args = append(args, "-p", f)
	}
	for _, f := range a.EnvFile {
		args = append(args, "--env-file", f)
	}
	for _, f := range a.Links {
		args = append(args, "--link", f)
	}

	args = append(args, fmt.Sprintf("%s:%s", a.Image, a.Tag))

	if len(a.Command) != 0 && len(extra) == 0 {
		args = append(args, a.Command...)
	} else {
		args = append(args, extra...)
	}

	return exec.Command(dockerExe, args...)
}
