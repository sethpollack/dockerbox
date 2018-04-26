package repo

import (
	"io/ioutil"
	"os"

	"github.com/imdario/mergo"
	"github.com/sethpollack/dockerbox/io"
	"github.com/sethpollack/dockerbox/registry"
	yaml "gopkg.in/yaml.v2"
)

const (
	repoFile = "$HOME/.dockerbox/repo.yaml"
)

type Repo struct {
	Applets map[string]Applet `yaml:"applets"`
}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) Init() error {
	return r.loadFile(os.ExpandEnv(repoFile), "file")
}

func (r *Repo) save() error {
	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.ExpandEnv(repoFile), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) loadFile(filename string, fileType string) error {
	b, err := io.ReadConfig(filename, fileType)
	if err != nil {
		return err
	}

	tmp := &Repo{}
	err = yaml.Unmarshal(b, tmp)
	if err != nil {
		return err
	}

	mergo.Merge(r, tmp)

	return nil
}

func (r *Repo) Update(reg *registry.Registry) error {
	for _, conf := range reg.Configs {
		err := r.loadFile(os.ExpandEnv(conf.Path), conf.Type)
		if err != nil {
			return err
		}
	}

	err := r.save()
	if err != nil {
		return err
	}

	return nil
}
