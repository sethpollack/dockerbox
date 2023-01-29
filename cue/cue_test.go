package cue

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sethpollack/dockerbox/applet"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type configs struct {
	path string
	data string
}

func TestNew(t *testing.T) {
	tt := []struct {
		name     string
		envs     map[string]string
		configs  []configs
		files    []string
		expected *applet.Root
		err      error
	}{
		{
			name: "compiles configs",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "test"
							image: "test"
						}
					`,
				},
			},
			files: []string{"/root/test.dbx.cue"},
			expected: &applet.Root{
				Applets: map[string]applet.Applet{
					"test": {
						Name:        "test",
						Image:       "test",
						Tag:         "latest",
						Interactive: true,
						RM:          true,
						TTY:         true,
					},
				},
			},
		},
		{
			name: "fails to compile invalid configs",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "test"
							image: "test"
							invalid: "invalid"
						}
					`,
				},
			},
			files: []string{"/root/test.dbx.cue"},
			err:   errors.New("failed to compile /root/test.dbx.cue: applets.test: field not allowed: invalid:\n    /root/test.dbx.cue:2:22\n    /root/test.dbx.cue:5:8\n    schema.cue:1:10\n"),
		},
		{
			name: "fails to decode invalid configs",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: []
					`,
				},
			},
			files: []string{"/root/test.dbx.cue"},
			err:   errors.New("failed to decode cue: applets: cannot use value [] (type list) as struct\n"),
		},
		{
			name: "unifies multiple configs",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "test"
							image: "test"
						}
					`,
				},
				{
					path: "/src/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "test"
							image: "test"
							entrypoint: "test"
						}
					`,
				},
			},
			files: []string{"/root/test.dbx.cue", "/src/test.dbx.cue"},
			expected: &applet.Root{
				Applets: map[string]applet.Applet{
					"test": {
						Name:        "test",
						Image:       "test",
						Entrypoint:  "test",
						Tag:         "latest",
						Interactive: true,
						RM:          true,
						TTY:         true,
					},
				},
			},
		},
		{
			name: "fails to unify invalid configs",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "test"
							image: "test"
						}
					`,
				},
				{
					path: "/src/test.dbx.cue",
					data: `
						applets: test: #Applet & {
							name: "foo"
							image: "test"
						}
					`,
				},
			},
			files: []string{"/root/test.dbx.cue", "/src/test.dbx.cue"},
			err:   errors.New("failed to unify cue: applets.test.name: conflicting values \"foo\" and \"test\":\n    /root/test.dbx.cue:3:14\n    /src/test.dbx.cue:3:14\n"),
		},
		{
			name: "adds environment variables",
			envs: map[string]string{
				"FOO": "bar",
			},
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: { Name: "\(environ.FOO)" }
					`,
				},
			},
			files: []string{"/root/test.dbx.cue"},
			expected: &applet.Root{
				Applets: map[string]applet.Applet{
					"test": {
						Name: "bar",
					},
				},
			},
		},
		{
			name: "validates required fields",
			configs: []configs{
				{
					path: "/root/test.dbx.cue",
					data: `
						applets: test: #Applet &{
							name: "test"
						}
					`,
				},
			},
			files: []string{"/root/test.dbx.cue"},
			err:   errors.New("failed to validate cue: applets.test.image: incomplete value string:\n    schema.cue:4:10\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			for _, config := range tc.configs {
				err := fs.MkdirAll(filepath.Dir(config.path), 0755)
				if err != nil {
					t.Fatal(err)
				}

				err = afero.WriteFile(fs, config.path, []byte(config.data), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			for k, v := range tc.envs {
				err := os.Setenv(k, v)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual, err := New(fs, tc.files)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
