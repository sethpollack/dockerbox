package repo

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"golang.org/x/crypto/ssh/terminal"
)

const dockerExe = "docker"

type Applet struct {
	Name       string `yaml:"name"`
	WorkDir    string `yaml:"work_dir"`
	Entrypoint string `yaml:"entrypoint"`
	Restart    string `yaml:"restart"`
	Network    string `yaml:"network"`
	EnvFilter  string `yaml:"env_filter"`
	Hostname   string `yaml:"hostname"`

	RM          bool `yaml:"rm"`
	TTY         bool `yaml:"tty"`
	Interactive bool `yaml:"interactive"`
	Privileged  bool `yaml:"privileged"`
	Detach      bool `yaml:"detach"`
	Kill        bool `yaml:"kill"`
	AllEnvs     bool `yaml:"all_envs"`

	DNS          []string `yaml:"dns"`
	DNSSearch    []string `yaml:"dns_search"`
	DNSOption    []string `yaml:"dns_option"`
	Env          []string `yaml:"environment"`
	Volumes      []string `yaml:"volumes"`
	Ports        []string `yaml:"ports"`
	EnvFile      []string `yaml:"env_file"`
	Dependencies []string `yaml:"dependencies"`
	BeforeHooks  []string `yaml:"before_hooks"`
	AfterHooks   []string `yaml:"after_hooks"`
	Links        []string `yaml:"links"`

	Image string `yaml:"image"`
	Tag   string `yaml:"image_tag"`

	Command []string `yaml:"command"`
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
	a.KillCmd().Run()
}

func isTTY() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
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

	raw.BeforeHooks = append(raw.Dependencies, raw.BeforeHooks...)
	raw.Dependencies = nil

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
		args = append(args, "--workdir", os.ExpandEnv(a.WorkDir))
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
	if a.Hostname != "" {
		args = append(args, "--hostname", a.Hostname)
	}

	if a.RM {
		args = append(args, "--rm")
	}
	if a.Privileged {
		args = append(args, "--privileged")
	}
	if a.Detach {
		args = append(args, "--detach")
	}
	if a.Interactive {
		args = append(args, "--interactive")
	}
	if isTTY() && a.TTY {
		args = append(args, "--tty")
	}
	if a.AllEnvs {
		for _, f := range os.Environ() {
			if matched, _ := regexp.MatchString(a.EnvFilter, f); matched {
				args = append(args, "-e", f)
			}
		}
	}
	for _, f := range a.DNS {
		args = append(args, "--dns", f)
	}
	for _, f := range a.DNSSearch {
		args = append(args, "--dns-search", f)
	}
	for _, f := range a.DNSOption {
		args = append(args, "--dns-option", f)
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

	if len(a.Command) != 0 {
		args = append(args, a.Command...)
	}

	if len(extra) != 0 {
		args = append(args, extra...)
	}

	return exec.Command(dockerExe, args...)
}
