package registry

import (
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/sethpollack/dockerbox/io"
	yaml "gopkg.in/yaml.v2"
)

const (
	regFile = "/registry.yaml"
)

type Registry struct {
	rootDir string
	Repos   []*Repo `yaml:"repos"`
}

type Repo struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Type string `yaml:"type"`
}

func New(rootDir string) (*Registry, error) {
	reg := &Registry{rootDir: rootDir}
	err := reg.load()
	if err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry) Add(name, path string) {
	found := false
	for _, repo := range r.Repos {
		if repo.Name == name {
			repo.Path = path
			repo.Type = getType(path)
			found = true
		}
	}

	if !found {
		r.Repos = append(r.Repos, &Repo{
			Name: name,
			Path: path,
			Type: getType(path),
		})
	}
}

func (r *Registry) Remove(name string) {
	for i, repo := range r.Repos {
		if repo.Name == name {
			r.Repos = append(r.Repos[:i], r.Repos[i+1:]...)
			break
		}
	}
}

func getType(path string) string {
	if isValidUrl(path) {
		return "url"
	} else {
		return "file"
	}
}

func (r *Registry) Save() error {
	b, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.rootDir+regFile, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) load() error {
	b, err := io.ReadFile(r.rootDir + regFile)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, r)
	if err != nil {
		return err
	}

	return nil
}

func isValidUrl(path string) bool {
	_, err := url.Parse(path)
	return (strings.Contains(path, "http://") || strings.Contains(path, "https://")) && err == nil
}
