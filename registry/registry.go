package registry

import (
	"github.com/sethpollack/dockerbox/io"
	yaml "gopkg.in/yaml.v2"
)

const (
	regFile = "/registry.yaml"
)

type Registry struct {
	Repos []*Repo `yaml:"repos"`
}

type Repo struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

func New(rootDir string) (*Registry, error) {
	reg := &Registry{}
	err := reg.load(rootDir + regFile)
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
