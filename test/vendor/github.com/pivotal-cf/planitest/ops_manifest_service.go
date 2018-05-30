package planitest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

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

func (o OpsManifestService) appendProperties(additionalProperties map[string]interface{}) ([]byte, error) {
	configInput, err := ioutil.ReadFile(o.config.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to configure product: %s (config file %s)", err, o.config.ConfigFile)
	}

	var inputJSON ConfigJSON
	err = json.Unmarshal(configInput, &inputJSON)

	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %s", err)
	}

	if inputJSON.NetworkConfig == nil {
		return nil, fmt.Errorf("network config must be provided in the config file")
	}

	if inputJSON.ProductProperties == nil {
		return nil, fmt.Errorf("product properties must be provided in the config file")
	}

	propertiesJSON, err := json.Marshal(inputJSON.ProductProperties)
	if err != nil {
		return nil, err // un-tested
	}

	var minimalProperties *map[string]interface{}
	err = json.Unmarshal(propertiesJSON, &minimalProperties)
	if err != nil {
		return nil, fmt.Errorf("could not parse product properties: %s", err)
	}

	combinedProperties := mergeProperties(*minimalProperties, additionalProperties)

	propertiesJSON, err = json.Marshal(combinedProperties)
	if err != nil {
		return nil, err // un-tested
	}

	err = json.Unmarshal(propertiesJSON, &inputJSON.ProductProperties)
	if err != nil {
		return nil, err // un-tested
	}

	outputJSON, err := json.Marshal(inputJSON)
	if err != nil {
		return nil, err
	}

	return outputJSON, nil
}

func (o OpsManifestService) RenderManifest(additionalProperties map[string]interface{}) (Manifest, error) {
	propertiesJSON, err := o.appendProperties(additionalProperties)
	if err != nil {
		panic(err)
	}

	manifest, err := o.opsManifestRunner.GetManifest(string(propertiesJSON), o.config.MetadataFile)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve bosh manifest: %s", err)
	}

	y, err := yaml.Marshal(manifest)
	if err != nil {
		return Manifest{}, err // un-tested
	}

	return NewManifest(string(y)), nil
}
