package applet

import (
	"fmt"
	"os"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

const dockerExe = "docker"

type Root struct {
	Ignore  map[string]Applet `json:"ignore"`
	Applets Applets           `json:"applets"`
}

type Applets map[string]Applet

type Applet struct {
	Entrypoint string `json:"entrypoint" flag:"entrypoint" desc:"Overwrite the default ENTRYPOINT of the image"`
	Hostname   string `json:"hostname" flag:"hostname" desc:"Container host name"`
	Image      string `json:"image" flag:"image" desc:"Container image"`
	Name       string `json:"name" flag:"name" desc:"Assign a name to the container"`
	Network    string `json:"network" flag:"network" desc:"Connect a container to a network"`
	Restart    string `json:"restart" flag:"restart" desc:"Restart policy to apply when a container exits (default \no\")"`
	Tag        string `json:"image_tag" flag:"tag" desc:"Container image tag"`
	WorkDir    string `json:"work_dir" flag:"workdir w" desc:"Working directory inside the container"`

	Detach      bool `json:"detach" flag:"detach d" desc:"Run container in background and print container ID"`
	Interactive bool `json:"interactive" flag:"interactive i" desc:"Keep STDIN open even if not attached"`
	Kill        bool `json:"kill" flag:"kill" desc:"Kill previous run on container with same name"`
	Privileged  bool `json:"privileged" flag:"privileged" desc:"Give extended privileges to this container"`
	Pull        bool `json:"pull" flag:"pull" desc:"Pull image before running it"`
	RM          bool `json:"rm" flag:"rm" desc:"Automatically remove the container when it exits"`
	TTY         bool `json:"tty" flag:"tty t" desc:"Allocate a pseudo-TTY"`

	AfterHooks  []Applet `json:"after_hooks" flag:"after-hook" desc:"Run container after"`
	BeforeHooks []Applet `json:"before_hooks" flag:"before-hook" desc:"Run container before."`
	Command     []string `json:"command" flag:"command" desc:"Command to run in container"`
	DNS         []string `json:"dns" flag:"dns" desc:"Set custom DNS servers"`
	DNSOption   []string `json:"dns_option" flag:"dns-option" desc:"Set DNS options"`
	DNSSearch   []string `json:"dns_search" flag:"dns-search" desc:"Set custom DNS search domains"`
	Env         []string `json:"environment" flag:"environment e" desc:"Set environment variables"`
	EnvFile     []string `json:"env_file" flag:"env-file" desc:"Read in a file of environment variables"`
	Links       []string `json:"links" flag:"link" desc:"Add link to another container"`
	Ports       []string `json:"ports" flag:"publish p" desc:"Publish a container's port(s) to the host"`
	Volumes     []string `json:"volumes" flag:"volume v" desc:"Bind mount a volume"`
}

func (root *Root) Compile(cfg *dockerbox.Config) ([][]string, error) {
	a, ok := root.Applets[cfg.EntryPoint]
	if !ok {
		return nil, fmt.Errorf("applet %s not found", cfg.EntryPoint)
	}

	fSet := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	err := gpflag.ParseTo(&a, fSet)
	if err != nil {
		return nil, fmt.Errorf("failed to create flag set from applet: %v", err)
	}

	dArgs, aArgs := splitArgs(cfg.Separator, cfg.Args)
	err = fSet.Parse(dArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse applet flags: %v", err)
	}

	err = root.validate(a)
	if err != nil {
		return nil, fmt.Errorf("failed to validate applet: %v", err)
	}

	return root.Applets.allCmds(a, aArgs...)
}

func (root *Root) validate(a Applet) error {
	err := a.validateRequired()
	if err != nil {
		return err
	}

	err = a.validateHooks(root.Applets)
	if err != nil {
		return err
	}

	return nil
}

func (applets Applets) allCmds(current Applet, args ...string) ([][]string, error) {
	var allCmds func(Applet, ...string) ([][]string, error)
	allCmds = func(applet Applet, args ...string) ([][]string, error) {
		cmds := [][]string{}

		for _, bh := range applet.BeforeHooks {
			h, ok := applets[bh.Name]
			if !ok {
				return nil, fmt.Errorf("before hook %s not found", bh.Name)
			}

			cmd, err := allCmds(h)
			if err != nil {
				return cmds, err
			}

			cmds = append(cmds, cmd...)
		}

		cmds = append(
			cmds,
			applet.appletCmds(args...)...,
		)

		for _, ah := range applet.AfterHooks {
			h, ok := applets[ah.Name]
			if !ok {
				return cmds, fmt.Errorf("after hook %s not found", ah.Name)
			}

			cmd, err := allCmds(h)
			if err != nil {
				return cmds, err
			}

			cmds = append(cmds, cmd...)
		}

		return cmds, nil
	}

	cmds, err := allCmds(current, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get commands: %v", err)
	}

	return cmds, nil
}

func (a Applet) killCmd() []string {
	args := []string{
		dockerExe,
		"kill",
		a.Name,
	}

	return args
}

func (a Applet) pullCmd() []string {
	args := []string{
		dockerExe,
		"pull",
	}

	if a.Tag != "" {
		args = append(args, fmt.Sprintf("%s:%s", a.Image, a.Tag))
	} else {
		args = append(args, a.Image)
	}

	return args
}

func (a Applet) runCmd(extra ...string) []string {
	args := []string{
		dockerExe,
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
		args = append(args, "-e", f)
	}

	for _, f := range a.Volumes {
		args = append(args, "-v", f)
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

	if a.Tag != "" {
		args = append(args, fmt.Sprintf("%s:%s", a.Image, a.Tag))
	} else {
		args = append(args, a.Image)
	}

	if len(a.Command) != 0 {
		args = append(args, a.Command...)
	}

	args = append(args, extra...)

	return args
}

func (a Applet) appletCmds(extra ...string) [][]string {
	commands := [][]string{}
	if a.Pull {
		commands = append(
			commands,
			a.pullCmd(),
		)
	}

	if a.Kill {
		commands = append(
			commands,
			a.killCmd(),
		)
	}

	commands = append(
		commands,
		a.runCmd(extra...),
	)

	return commands
}

func (a Applet) validateRequired() error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}

	if a.Image == "" {
		return fmt.Errorf("image is required")
	}

	return nil
}

func (a Applet) validateHooks(applets Applets) error {
	var validate func(Applet, map[string]bool) error

	validate = func(applet Applet, visited map[string]bool) error {
		if visited[applet.Name] {
			return fmt.Errorf("circular dependency detected: %s", applet.Name)
		}

		visited[applet.Name] = true

		for _, h := range applet.BeforeHooks {
			bh, ok := applets[h.Name]
			if !ok {
				return fmt.Errorf("before hook %s not found", h.Name)
			}

			err := validate(bh, visited)
			if err != nil {
				return err
			}
		}

		for _, h := range applet.AfterHooks {
			ah, ok := applets[h.Name]
			if !ok {
				return fmt.Errorf("after hook %s not found", h.Name)
			}

			err := validate(ah, visited)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return validate(a, map[string]bool{})
}

func isTTY() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}

func splitArgs(separator string, args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == separator {
			return args[:i], args[i+1:]
		}
	}
	return []string{}, args
}
