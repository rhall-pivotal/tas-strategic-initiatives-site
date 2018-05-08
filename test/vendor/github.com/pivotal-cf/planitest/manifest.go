package planitest

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

//go:generate counterfeiter -o ./fakes/command_runner.go --fake-name CommandRunner . CommandRunner
type CommandRunner interface {
	Run(string, ...string) (string, string, error)
}

type Manifest struct {
	cmdRunner CommandRunner
	content   string
}

func NewManifest(content string) Manifest {
	return Manifest{
		cmdRunner: NewExecutor(),
		content:   content,
	}
}

func NewManifestWithRunner(content string, cmdRunner CommandRunner) Manifest {
	return Manifest{
		cmdRunner: cmdRunner,
		content:   content,
	}
}

func (m *Manifest) FindInstanceGroupJob(instanceGroup, job string) (Manifest, error) {
	path := fmt.Sprintf("/instance_groups/name=%s/jobs/name=%s", instanceGroup, job)

	output, errOutput, err := m.interpolate(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("Problem interpolating manifest: %s: %s", err, errOutput)
	}
	return NewManifestWithRunner(output, m.cmdRunner), nil
}

func (m *Manifest) Property(path string) (interface{}, error) {
	resultYAML, errOutput, err := m.interpolate(fmt.Sprintf("/properties/%s", path))
	if err != nil {
		return "", fmt.Errorf("Problem interpolating manifest: %s: %s", err, errOutput)
	}

	var result interface{}
	err = yaml.Unmarshal([]byte(resultYAML), &result)
	if err != nil {
		return nil, fmt.Errorf("Problem unmarshalling result retrieving property %q: %s", path, err)
	}
	return result, nil
}

func (m *Manifest) String() string {
	return m.content
}

func (m *Manifest) interpolate(expr string) (string, string, error) {
	manifestFile, err := ioutil.TempFile("", "manifest")
	if err != nil {
		return "", "", err // un-tested
	}

	defer os.Remove(manifestFile.Name()) // un-tested

	if _, err = manifestFile.WriteString(m.content); err != nil {
		return "", "", err // un-tested
	}

	if err = manifestFile.Close(); err != nil {
		return "", "", err // un-tested
	}

	return m.cmdRunner.Run(
		"bosh",
		"--non-interactive",
		"interpolate",
		fmt.Sprintf("--path=%s", expr),
		manifestFile.Name(),
	)
}
