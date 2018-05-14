package planitest

import (
	"fmt"

	"github.com/cppforlife/go-patch/patch"
	yaml "gopkg.in/yaml.v2"
)

//go:generate counterfeiter -o ./fakes/command_runner.go --fake-name CommandRunner . CommandRunner
type CommandRunner interface {
	Run(string, ...string) (string, string, error)
}

type Manifest struct {
	content string
}

func NewManifest(content string) Manifest {
	return Manifest{
		content: content,
	}
}

func (m *Manifest) FindInstanceGroupJob(instanceGroup, job string) (Manifest, error) {
	path := fmt.Sprintf("/instance_groups/name=%s/jobs/name=%s", instanceGroup, job)

	result, err := m.interpolate(path)
	if err != nil {
		return Manifest{}, err
	}

	content, err := yaml.Marshal(result)
	if err != nil {
		return Manifest{}, err // should never happen
	}
	return NewManifest(string(content)), nil
}

func (m *Manifest) Property(path string) (interface{}, error) {
	return m.interpolate(fmt.Sprintf("/properties/%s", path))
}

func (m *Manifest) Path(path string) (interface{}, error) {
	return m.interpolate(path)
}

func (m *Manifest) String() string {
	return m.content
}

func (m *Manifest) interpolate(path string) (interface{}, error) {
	var content interface{}
	err := yaml.Unmarshal([]byte(m.content), &content)
	if err != nil {
		return "", fmt.Errorf("failed to parse manifest: %s", err)
	}

	res, err := patch.FindOp{Path: patch.MustNewPointerFromString(path)}.Apply(content)
	if err != nil {
		return "", fmt.Errorf("failed to find value at path '%s': %s", path, m.content)
	}

	return res, nil
}
