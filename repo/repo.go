package repo

import (
	"io/ioutil"
	"os"

	"github.com/sethpollack/dockerbox/io"
	"github.com/sethpollack/dockerbox/registry"
	yaml "gopkg.in/yaml.v2"
)

const (
	cacheFile = "/.cache.yaml"
)

type Repo struct {
	rootDir string
	Applets map[string]Applet `yaml:"applets"`
}

func New(rootDir string) *Repo {
	return &Repo{rootDir: rootDir}
}

func (r *Repo) Init() error {
	return r.loadFile(r.rootDir+cacheFile, "file")
}

func (r *Repo) save() error {
	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.rootDir+cacheFile, b, 0644)
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

	err = yaml.Unmarshal(b, r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Update(reg *registry.Registry) error {
	for _, repo := range reg.Repos {
		err := r.loadFile(os.ExpandEnv(repo.Path), repo.Type)
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
