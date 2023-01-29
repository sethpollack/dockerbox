package dockerbox

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tt := []struct {
		name string
		envs map[string]string
		cfg  *Config
		err  error
	}{
		{
			name: "defaults",
			envs: map[string]string{
				"HOME": "/root",
			},
			cfg: &Config{
				RootDir:      "/root/.dockerbox",
				InstallDir:   "/root/.dockerbox/bin",
				Separator:    "--",
				WD:           "",
				DockerboxExe: "",
				EntryPoint:   "",
				Args:         []string{},
			},
		},
		{
			name: "env var overrides",
			envs: map[string]string{
				"DOCKERBOX_ROOT_DIR":    "/foo",
				"DOCKERBOX_INSTALL_DIR": "/foo/bin",
				"DOCKERBOX_SEPARATOR":   "***",
			},
			cfg: &Config{
				RootDir:      "/foo",
				InstallDir:   "/foo/bin",
				Separator:    "***",
				WD:           "",
				DockerboxExe: "",
				EntryPoint:   "",
				Args:         []string{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Unsetenv("HOME")
			os.Unsetenv("DOCKERBOX_ROOT_DIR")
			os.Unsetenv("DOCKERBOX_INSTALL_DIR")
			os.Unsetenv("DOCKERBOX_SEPARATOR")

			for k, v := range tc.envs {
				os.Setenv(k, v)
			}

			cfg, err := New("", "", "", []string{})

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.cfg, cfg)
		})
	}
}

func TestGetConfigurations(t *testing.T) {
	tt := []struct {
		name string

		configs  []string
		expected []string

		wd   string
		root string

		err error
	}{
		{
			name: "finds files relative to the working directory",
			wd:   "/src/foo",
			root: "/root",
			configs: []string{
				"/src/foo/test.dbx.cue",
				"/src/test.dbx.cue",
			},
			expected: []string{
				"/src/foo/test.dbx.cue",
				"/src/test.dbx.cue",
			},
		},
		{
			name: "only find files with the .dbx.cue extension",
			wd:   "/src/foo",
			root: "/root",
			configs: []string{
				"/src/foo/test.dbx.cue",
				"/src/test.dbx.cue",
				"/src/test.cue",
			},
			expected: []string{
				"/src/foo/test.dbx.cue",
				"/src/test.dbx.cue",
			},
		},
		{
			name: "finds files in the root directory",
			wd:   "/src/foo",
			root: "/root",
			configs: []string{
				"/root/test.dbx.cue",
			},
			expected: []string{
				"/root/test.dbx.cue",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			err := fs.MkdirAll(tc.root, 0755)
			if err != nil {
				t.Fatal(err)
			}

			err = fs.MkdirAll(tc.wd, 0755)
			if err != nil {
				t.Fatal(err)
			}

			for _, config := range tc.configs {
				err = afero.WriteFile(fs, config, []byte(""), 0644)
				if err != nil {
					t.Fatal(err)
				}
			}

			actual, err := GetConfigurations(fs, tc.wd, tc.root)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
