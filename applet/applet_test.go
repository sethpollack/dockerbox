package applet

import (
	"errors"
	"testing"

	"github.com/sethpollack/dockerbox/dockerbox"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	tt := []struct {
		name string
		root Root
		cmds [][]string
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
			cmds: [][]string{
				{"docker", "pull", "test"},
				{"docker", "run", "test"},
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
			cmds: [][]string{
				{"docker", "kill", "test"},
				{"docker", "run", "--name", "test", "test"},
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
			cmds: [][]string{
				{"docker", "run", "before"},
				{"docker", "run", "test"},
				{"docker", "run", "after"},
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
						Network:    "test",
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
			cmds: [][]string{
				{"docker", "run", "--name", "before", "before"},
				{"docker", "pull", "test:test"},
				{"docker", "kill", "test"},
				{"docker", "run", "--name", "test", "--workdir", "test", "--entrypoint", "test", "--restart", "test", "--network", "test", "--hostname", "test", "--rm", "--privileged", "--detach", "--interactive", "--dns", "test", "--dns-search", "test", "--dns-option", "test", "-e", "test", "-v", "test", "-p", "test", "--env-file", "test", "--link", "test", "test:test", "test", "my", "args"},
				{"docker", "run", "--name", "after", "after"},
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
