package cue

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"github.com/sethpollack/dockerbox/applet"
	"github.com/spf13/afero"
)

//go:embed schema.cue
var schema []byte

type Cue struct {
	ctx   *cue.Context
	fs    afero.Fs
	files []string
}

func New(fs afero.Fs, files []string) (*applet.Root, error) {
	c := &Cue{
		fs:    fs,
		files: files,
		ctx:   cuecontext.New(),
	}

	return c.Compile()
}

func (c *Cue) Compile() (*applet.Root, error) {
	values, err := c.Values()
	if err != nil {
		return nil, err
	}

	value := c.Unify(values)
	if value.Err() != nil {
		return nil, fmt.Errorf("failed to unify cue: %s", errors.Details(value.Err(), nil))
	}

	err = value.Validate(
		cue.Final(),
		cue.Concrete(true),
		cue.DisallowCycles(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to validate cue: %s", errors.Details(err, nil))
	}

	root := &applet.Root{}

	if err := value.Decode(root); err != nil {
		return nil, fmt.Errorf("failed to decode cue: %v", errors.Details(err, nil))
	}

	return root, nil
}

func (c *Cue) Unify(values []cue.Value) cue.Value {
	value := values[0]

	for _, v := range values[1:] {
		value = value.Unify(v)
	}

	return value
}

func (c *Cue) AddEnvs(v cue.Value) cue.Value {
	envs := map[string]string{}
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		envs[pair[0]] = pair[1]
	}

	return v.Unify(
		c.ctx.Encode(map[string]any{
			"environ": envs,
		}),
	)
}

func (c *Cue) CompileSchema() cue.Value {
	value := c.ctx.CompileBytes(
		schema,
		cue.Filename("schema.cue"),
	)

	return c.AddEnvs(value)
}

func (c *Cue) Values() ([]cue.Value, error) {
	values := []cue.Value{}

	schema := c.CompileSchema()
	if schema.Err() != nil {
		return nil, fmt.Errorf("failed to compile schema: %v", errors.Details(schema.Err(), nil))
	}

	for _, filename := range c.files {
		bytes, err := afero.ReadFile(c.fs, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", filename, err)
		}
		value := c.ctx.CompileBytes(
			bytes,
			cue.Filename(filename),
			cue.Scope(schema),
		)
		if value.Err() != nil {
			return nil, fmt.Errorf("failed to compile %s: %v", filename, errors.Details(value.Err(), nil))
		}
		values = append(values, value)
	}

	return values, nil
}
