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
	Name       string `yaml:"name" flag:"name" desc:"Assign a name to the container"`
	WorkDir    string `yaml:"work_dir" flag:"workdir w" desc:"Working directory inside the container"`
	Entrypoint string `yaml:"entrypoint" flag:"entrypoint" desc:"Overwrite the default ENTRYPOINT of the image"`
	Restart    string `yaml:"restart" flag:"restart" desc:"Restart policy to apply when a container exits (default \no\")"`
	Network    string `yaml:"network" flag:"network" desc:"Connect a container to a network"`
	EnvFilter  string `yaml:"env_filter" flag:"env-filter" desc:"Filter env vars passed to container from --all-envs"`
	Hostname   string `yaml:"hostname" flag:"hostname" desc:"Container host name"`
	Image      string `yaml:"image" flag:"image" desc:"Container image"`
	Tag        string `yaml:"image_tag" flag:"tag" desc:"Container image tag"`

	RM               bool `yaml:"rm" flag:"rm" desc:"Automatically remove the container when it exits"`
	TTY              bool `yaml:"tty" flag:"tty t" desc:"Allocate a pseudo-TTY"`
	Interactive      bool `yaml:"interactive" flag:"interactive i" desc:"Keep STDIN open even if not attached"`
	Privileged       bool `yaml:"privileged" flag:"privileged" desc:"Give extended privileges to this container"`
	Detach           bool `yaml:"detach" flag:"detach d" desc:"Run container in background and print container ID"`
	Kill             bool `yaml:"kill" flag:"kill" desc:"Kill previous run on container with same name"`
	AllEnvs          bool `yaml:"all_envs" flag:"all-envs" desc:"Pass all envars to container"`
	Pull             bool `yaml:"pull" flag:"pull" desc:"Pull image before running it"`
	InverseEnvFilter bool `yaml:"inverse" flag:"inverse" desc:"Inverse env-filter"`

	DNS          []string `yaml:"dns" flag:"dns" desc:"Set custom DNS servers"`
	DNSSearch    []string `yaml:"dns_search" flag:"dns-search" desc:"Set custom DNS search domains"`
	DNSOption    []string `yaml:"dns_option" flag:"dns-option" desc:"Set DNS options"`
	Env          []string `yaml:"environment" flag:"environment e" desc:"Set environment variables"`
	Volumes      []string `yaml:"volumes" flag:"volume v" desc:"Bind mount a volume"`
	Ports        []string `yaml:"ports" flag:"publish p" desc:"Publish a container's port(s) to the host"`
	EnvFile      []string `yaml:"env_file" flag:"env-file" desc:"Read in a file of environment variables"`
	Dependencies []string `yaml:"dependencies" flag:"dependency" desc:"Run container before"`
	BeforeHooks  []string `yaml:"before_hooks" flag:"before-hook" desc:"Run container before."`
	AfterHooks   []string `yaml:"after_hooks" flag:"after-hook" desc:"Run container after"`
	Links        []string `yaml:"links" flag:"link" desc:"Add link to another container"`
	Command      []string `yaml:"command" flag:"command" desc:"Command to run in container"`
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
	if a.Kill {
		a.KillCmd().Run()
	}
	if a.Pull {
		a.PullCmd().Run()
	}
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

func (a *Applet) PullCmd() *exec.Cmd {
	args := []string{
		"pull",
		fmt.Sprintf("%s:%s", a.Image, a.Tag),
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
			if matched, _ := regexp.MatchString(a.EnvFilter, f); a.InverseEnvFilter != matched {
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
