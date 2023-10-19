package applet

import (
	"errors"
	"testing"

	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/sethpollack/dockerbox/runner"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	tt := []struct {
		name string
		root Root
		cmds []runner.Cmd
		cfg  *dockerbox.Config
		err  error
	}{
		{
			name: "no applet",
			root: Root{},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("applet test not found"),
		},
		{
			name: "validates missing name",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						Image: "test",
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: applet_name is required"),
		},
		{
			name: "validates missing image",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: image is required"),
		},
		{
			name: "validates missing before hook",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
						Image:      "test",
						BeforeHooks: []Applet{
							{
								AppletName: "before",
								Image:      "before",
							},
						},
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: before hook before not found"),
		},
		{
			name: "validates missing after hook",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
						Image:      "test",
						AfterHooks: []Applet{
							{
								AppletName: "after",
								Image:      "after",
							},
						},
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: after hook after not found"),
		},
		{
			name: "validates circular before hooks",
			root: Root{
				Applets: map[string]Applet{
					"before": {
						AppletName: "before",
						Image:      "before",
						BeforeHooks: []Applet{
							{
								AppletName: "test",
								Image:      "test",
							},
						},
					},
					"test": {
						AppletName: "test",
						Image:      "test",
						BeforeHooks: []Applet{
							{
								AppletName: "before",
								Image:      "before",
							},
						},
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: circular dependency detected: test"),
		},
		{
			name: "validates circular after hooks",
			root: Root{
				Applets: map[string]Applet{
					"after": {
						AppletName: "after",
						Image:      "after",
						AfterHooks: []Applet{
							{
								AppletName: "test",
								Image:      "test",
							},
						},
					},
					"test": {
						AppletName: "test",
						Image:      "test",
						AfterHooks: []Applet{
							{
								AppletName: "after",
								Image:      "after",
							},
						},
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: errors.New("failed to validate applet: circular dependency detected: test"),
		},
		{
			name: "pull arg",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
						Image:      "test",
						Pull:       true,
					},
				},
			},
			cmds: []runner.Cmd{
				{Args: []string{"docker", "pull", "test"}},
				{Args: []string{"docker", "run", "test"}},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: nil,
		},
		{
			name: "kill arg",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
						Name:       "test",
						Image:      "test",
						Kill:       true,
					},
				},
			},
			cmds: []runner.Cmd{
				{Silent: true, Args: []string{"docker", "kill", "test"}},
				{Args: []string{"docker", "run", "--name", "test", "test"}},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: nil,
		},
		{
			name: "hooks",
			root: Root{
				Applets: map[string]Applet{
					"before": {
						AppletName: "before",
						Image:      "before",
					},
					"after": {
						AppletName: "after",
						Image:      "after",
					},
					"test": {
						AppletName: "test",
						Image:      "test",
						BeforeHooks: []Applet{
							{
								AppletName: "before",
								Image:      "before",
							},
						},
						AfterHooks: []Applet{
							{
								AppletName: "after",
								Image:      "after",
							},
						},
					},
				},
			},
			cmds: []runner.Cmd{
				{Args: []string{"docker", "run", "before"}},
				{Args: []string{"docker", "run", "test"}},
				{Args: []string{"docker", "run", "after"}},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
			},
			err: nil,
		},
		{
			name: "invalid flags",
			root: Root{
				Applets: map[string]Applet{
					"test": {
						AppletName: "test",
						Image:      "test",
					},
				},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
				Args:       []string{"--invalid", "---", "my", "args"},
				Separator:  "---",
			},
			err: errors.New("failed to parse applet flags: unknown flag: --invalid"),
		},
		{
			name: "full applet",
			root: Root{
				Volumes: map[string]Volume{
					"test": {
						Name:   "test",
						Driver: "test",
					},
				},
				Networks: map[string]Network{
					"test": {
						Name:   "test",
						Driver: "test",
					},
				},
				Applets: map[string]Applet{
					"before": {
						AppletName: "before",
						Name:       "before",
						Image:      "before",
					},
					"after": {
						AppletName: "after",
						Name:       "after",
						Image:      "after",
					},
					"test": {
						AppletName: "test",
						Name:       "test",
						Image:      "test",
						Tag:        "test",
						Entrypoint: "test",
						Hostname:   "test",
						Restart:    "test",
						WorkDir:    "test",

						Detach:      true,
						Interactive: true,
						TTY:         true,
						Privileged:  true,
						Pull:        true,
						Kill:        true,
						RM:          true,

						Command:   []string{"test"},
						DNS:       []string{"test"},
						DNSOption: []string{"test"},
						DNSSearch: []string{"test"},
						Env:       []string{"test"},
						EnvFile:   []string{"test"},
						Links:     []string{"test"},
						Volumes:   []string{"test"},
						Networks:  []string{"test"},
						Ports:     []string{"test"},

						BeforeHooks: []Applet{
							{
								AppletName: "before",
								Name:       "before",
								Image:      "before",
							},
						},
						AfterHooks: []Applet{
							{
								AppletName: "after",
								Name:       "after",
								Image:      "after",
							},
						},
					},
				},
			},
			cmds: []runner.Cmd{
				{Args: []string{"docker", "network", "create", "--driver", "test", "test"}},
				{Args: []string{"docker", "volume", "create", "--driver", "test", "test"}},
				{Args: []string{"docker", "run", "--name", "before", "before"}},
				{Args: []string{"docker", "pull", "test:test"}},
				{Silent: true, Args: []string{"docker", "kill", "test"}},
				{Args: []string{"docker", "run", "--name", "test", "--workdir", "test", "--entrypoint", "test", "--restart", "test", "--hostname", "test", "--rm", "--privileged", "--detach", "--interactive", "--dns", "test", "--dns-search", "test", "--dns-option", "test", "-e", "test", "-v", "test", "--network", "test", "-p", "test", "--env-file", "test", "--link", "test", "test:test", "test", "my", "args"}},
				{Args: []string{"docker", "run", "--name", "after", "after"}},
			},
			cfg: &dockerbox.Config{
				EntryPoint: "test",
				Args:       []string{"--tty", "---", "my", "args"},
				Separator:  "---",
			},
			err: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.root.Compile(tc.cfg)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.cmds, actual)
		})
	}
}
