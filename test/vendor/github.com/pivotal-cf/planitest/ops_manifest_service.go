package planitest

import (
	"fmt"

	"github.com/pivotal-cf/planitest/internal"
	yaml "gopkg.in/yaml.v2"
)

type OpsManifestService struct {
	config            OpsManifestConfig
	opsManifestRunner OpsManifestRunner
}

type OpsManifestConfig struct {
	ConfigFile   string
	MetadataFile string
}

//go:generate counterfeiter -o ./fakes/ops_manifest_runner.go --fake-name OpsManifestRunner . OpsManifestRunner
type OpsManifestRunner interface {
	GetManifest(productProperties, metadataFilePath string) (map[string]interface{}, error)
}

func NewOpsManifestServiceWithRunner(config OpsManifestConfig, opsManifestRunner OpsManifestRunner) (*OpsManifestService, error) {
	return &OpsManifestService{config: config, opsManifestRunner: opsManifestRunner}, nil
}

func NewOpsManifestService(config OpsManifestConfig) (*OpsManifestService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	opsManifestRunner := internal.NewOpsManifestRunner(NewExecutor())
	return NewOpsManifestServiceWithRunner(config, opsManifestRunner)
}

func (o OpsManifestService) RenderManifest(additionalProperties map[string]interface{}) (Manifest, error) {
	configInput, err := MergeAdditionalProductProperties(o.config.ConfigFile, additionalProperties)
	if err != nil {
		return Manifest{}, err
	}

	manifest, err := o.opsManifestRunner.GetManifest(string(configInput), o.config.MetadataFile)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve bosh manifest: %s", err)
	}

	y, err := yaml.Marshal(manifest)
	if err != nil {
		return Manifest{}, err // un-tested
	}

	return NewManifest(string(y)), nil
}
