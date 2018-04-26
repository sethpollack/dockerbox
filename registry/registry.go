package registry

import (
	"os"

	"github.com/sethpollack/dockerbox/io"
	yaml "gopkg.in/yaml.v2"
)

const (
	regFile = "$HOME/.dockerbox/registry.yaml"
)

type Registry struct {
	Configs []*Config `yaml:"configs"`
}

type Config struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

func New() (*Registry, error) {
	reg := &Registry{}
	err := reg.load(os.ExpandEnv(regFile))
	if err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry) load(filename string) error {
	b, err := io.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, r)
	if err != nil {
		return err
	}

	return nil
}
