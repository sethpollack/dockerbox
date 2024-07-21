package applet

import (
	"fmt"
	"os"

	"github.com/octago/sflags/gen/gpflag"
	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/sethpollack/dockerbox/runner"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

const dockerExe = "docker"

type Root struct {
	Ignore   map[string]Applet  `json:"ignore"`
	Applets  Applets            `json:"applets"`
	Volumes  map[string]Volume  `json:"volumes"`
	Networks map[string]Network `json:"networks"`
}

type Applets map[string]Applet

type Applet struct {
	AppletName string `json:"applet_name" desc:"name of the applet"`

	Entrypoint string `json:"entrypoint" flag:"entrypoint" desc:"Overwrite the default ENTRYPOINT of the image"`
	Hostname   string `json:"hostname" flag:"hostname" desc:"Container host name"`
	Image      string `json:"image" flag:"image" desc:"Container image"`
	Name       string `json:"name" flag:"name" desc:"Assign a name to the container"`
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
	Networks    []string `json:"network" flag:"network" desc:"Connect a container to a network"`
	Volumes     []string `json:"volumes" flag:"volume v" desc:"Bind mount a volume"`
}

type Volume struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

type Network struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

func (root *Root) Compile(cfg *dockerbox.Config) ([]runner.Cmd, error) {
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

	allCmds := []runner.Cmd{}

	for _, n := range root.Networks {
		allCmds = append(allCmds, n.createNetworkCmd())
	}

	for _, v := range a.Volumes {
		if vol, ok := root.Volumes[v]; ok {
			allCmds = append(allCmds, vol.createVolumeCmd())
		}
	}

	appletCmds, err := root.Applets.allCmds(a, aArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to get applet commands: %v", err)
	}

	allCmds = append(allCmds, appletCmds...)

	return allCmds, nil
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

func (applets Applets) allCmds(current Applet, args ...string) ([]runner.Cmd, error) {
	var allCmds func(Applet, ...string) ([]runner.Cmd, error)
	allCmds = func(applet Applet, args ...string) ([]runner.Cmd, error) {
		cmds := []runner.Cmd{}

		for _, bh := range applet.BeforeHooks {
			h, ok := applets[bh.AppletName]
			if !ok {
				return nil, fmt.Errorf("before hook %s not found", bh.AppletName)
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
			h, ok := applets[ah.AppletName]
			if !ok {
				return cmds, fmt.Errorf("after hook %s not found", ah.AppletName)
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

func (a Applet) killCmd() runner.Cmd {
	return runner.Cmd{
		Silent: true,
		Args: []string{
			dockerExe,
			"kill",
			a.Name,
		},
	}
}

func (a Applet) pullCmd() runner.Cmd {
	args := []string{
		dockerExe,
		"pull",
	}

	if a.Tag != "" {
		args = append(args, fmt.Sprintf("%s:%s", a.Image, a.Tag))
	} else {
		args = append(args, a.Image)
	}

	return runner.Cmd{
		Args: args,
	}
}

func (a Applet) runCmd(extra ...string) runner.Cmd {
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

	for _, f := range a.Networks {
		args = append(args, "--network", f)
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

	return runner.Cmd{
		Args: args,
	}
}

func (a Applet) appletCmds(extra ...string) []runner.Cmd {
	commands := []runner.Cmd{}
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
	if a.AppletName == "" {
		return fmt.Errorf("applet_name is required")
	}

	if a.Image == "" {
		return fmt.Errorf("image is required")
	}

	return nil
}

func (a Applet) validateHooks(applets Applets) error {
	var validate func(Applet, map[string]bool) error

	validate = func(applet Applet, visited map[string]bool) error {
		if visited[applet.AppletName] {
			return fmt.Errorf("circular dependency detected: %s", applet.AppletName)
		}

		visited[applet.AppletName] = true

		for _, h := range applet.BeforeHooks {
			bh, ok := applets[h.AppletName]
			if !ok {
				return fmt.Errorf("before hook %s not found", h.AppletName)
			}

			err := validate(bh, visited)
			if err != nil {
				return err
			}
		}

		for _, h := range applet.AfterHooks {
			ah, ok := applets[h.AppletName]
			if !ok {
				return fmt.Errorf("after hook %s not found", h.AppletName)
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

func (v Volume) createVolumeCmd() runner.Cmd {
	args := []string{
		dockerExe,
		"volume",
		"create",
	}

	if v.Driver != "" {
		args = append(args, "--driver", v.Driver)
	}

	args = append(args, v.Name)

	return runner.Cmd{
		Args: args,
	}
}

func (n Network) createNetworkCmd() runner.Cmd {
	args := []string{
		dockerExe,
		"network",
		"create",
	}

	if n.Driver != "" {
		args = append(args, "--driver", n.Driver)
	}

	args = append(args, n.Name)

	return runner.Cmd{
		Args: args,
	}
}

func isTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

func splitArgs(separator string, args []string) ([]string, []string) {
	for i, arg := range args {
		if arg == separator {
			return args[:i], args[i+1:]
		}
	}
	return []string{}, args
}
