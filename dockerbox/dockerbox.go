package dockerbox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/afero"
)

type Config struct {
	RootDir    string `envconfig:"DOCKERBOX_ROOT_DIR" default:"$HOME/.dockerbox"`
	InstallDir string `envconfig:"DOCKERBOX_INSTALL_DIR" default:"$HOME/.dockerbox/bin"`
	Separator  string `envconfig:"DOCKERBOX_SEPARATOR" default:"--"`

	WD           string
	DockerboxExe string
	EntryPoint   string
	Args         []string
}

func New(ent, wd, exe string, args []string) (*Config, error) {
	cfg := &Config{
		EntryPoint:   ent,
		WD:           wd,
		DockerboxExe: exe,
		Args:         args,
	}

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	cfg.RootDir = os.ExpandEnv(cfg.RootDir)
	cfg.InstallDir = os.ExpandEnv(cfg.InstallDir)

	return cfg, nil
}

func GetConfigurations(fs afero.Fs, wd, rootDir string) ([]string, error) {
	files := []string{}

	dirFiles, err := readDir(fs, rootDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %v", rootDir, err)
	}
	files = append(files, dirFiles...)

	currentDir := wd
	for currentDir != "/" {
		dirFiles, err := readDir(fs, currentDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read dir %s: %v", currentDir, err)
		}
		files = append(files, dirFiles...)

		currentDir = filepath.Dir(currentDir)
	}

	return files, nil
}

func readDir(fs afero.Fs, currentDir string) ([]string, error) {
	files := []string{}

	dir, err := afero.ReadDir(fs, currentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %v", dir, err)
	}

	for _, file := range dir {
		ok, _ := filepath.Match("*.dbx.cue", file.Name())
		if !file.IsDir() && ok {
			files = append(files, filepath.Join(currentDir, file.Name()))
		}
	}

	return files, nil
}
