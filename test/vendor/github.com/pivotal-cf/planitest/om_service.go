package planitest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pivotal-cf/planitest/internal"
	yaml "gopkg.in/yaml.v2"
)

type OMService struct {
	omRunner OMRunner
	config   OMConfig
}

type OMConfig struct {
	Name       string
	Version    string
	ConfigFile string
}

//go:generate counterfeiter -o ./fakes/om_runner.go --fake-name OMRunner . OMRunner
type OMRunner interface {
	ResetAndConfigure(productName string, productVersion string, configJSON string) error
	GetManifest(productGUID string) (map[string]interface{}, error)
	FindStagedProduct(productName string) (internal.StagedProduct, error)
}

//go:generate counterfeiter -o ./fakes/file_io.go --fake-name FileIO . FileIO
type FileIO interface {
	TempFile(string, string) (*os.File, error)
	Remove(string) error
}

func NewOMService(config OMConfig) (*OMService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	omRunner := internal.NewOMRunner(NewExecutor())
	return NewOMServiceWithRunner(config, omRunner)
}

func NewOMServiceWithRunner(config OMConfig, omRunner OMRunner) (*OMService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	return &OMService{config: config, omRunner: omRunner}, nil
}

func TempFile(a, b string) (*os.File, error) {
	return ioutil.TempFile(a, b)
}

func Remove(a string) error {
	return os.Remove(a)
}

func validateEnvironmentVariables() error {
	requiredEnvVars := []string{"OM_USERNAME", "OM_PASSWORD", "OM_URL"}
	for _, envVar := range requiredEnvVars {
		value := os.Getenv(envVar)
		if value == "" {
			return fmt.Errorf("Environment variable %s must be set", envVar)
		}
	}
	return nil
}

func (o *OMService) appendProperties(additionalProperties map[string]interface{}) ([]byte, error) {
	configInput, err := ioutil.ReadFile(o.config.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to configure product %q: %s (config file %s)", o.config.Name, err, o.config.ConfigFile)
	}

	var inputJSON ConfigJSON
	err = json.Unmarshal(configInput, &inputJSON)

	if err != nil {
		return nil, fmt.Errorf("could not parse config file: %s", err)
	}

	if inputJSON.NetworkProperties == nil {
		return nil, fmt.Errorf("network properties must be provided in the config file")
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

	return outputJSON, nil
}

func mergeProperties(minimalProperties, additionalProperties map[string]interface{}) map[string]interface{} {
	combinedProperties := make(map[string]interface{}, len(minimalProperties)+len(additionalProperties))
	for k, v := range minimalProperties {
		combinedProperties[k] = v
	}

	for k, v := range additionalProperties {
		combinedProperties[k] = map[string]interface{}{
			"value": v,
		}
	}
	return combinedProperties
}

func (o OMService) RenderManifest(additionalProperties map[string]interface{}) (Manifest, error) {
	stagedProduct, err := o.omRunner.FindStagedProduct(o.config.Name)
	if err != nil {
		return Manifest{}, err
	}

	err = o.configure(additionalProperties)
	if err != nil {
		return Manifest{}, err
	}

	manifest, err := o.omRunner.GetManifest(stagedProduct.GUID)
	if err != nil {
		return Manifest{}, err
	}

	y, err := yaml.Marshal(manifest)
	if err != nil {
		return Manifest{}, err // un-tested
	}

	return NewManifest(string(y)), nil
}

func (o *OMService) configure(additionalProperties map[string]interface{}) error {
	configInput, err := o.appendProperties(additionalProperties)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s", o.config.Name, err)
	}

	err = o.omRunner.ResetAndConfigure(o.config.Name, o.config.Version, string(configInput))
	if err != nil {
		panic(err)
	}

	return nil
}
