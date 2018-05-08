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
	ResetAndConfigure(productName string, productVersion string, propertiesJSON string, networkJSON string) error
	GetManifest(productGUID string) (map[string]interface{}, error)
	FindStagedProduct(productName string) (internal.StagedProduct, error)
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

func extractPropertiesAndNetworkConfig(configInput []byte, additionalProperties map[string]interface{}) ([]byte, []byte, error) {
	var configJSON ConfigJSON
	err := json.Unmarshal(configInput, &configJSON)

	if err != nil {
		return nil, nil, fmt.Errorf("could not parse config file: %s", err)
	}

	if configJSON.NetworkConfig == nil {
		return nil, nil, fmt.Errorf("network config must be provided in the config file")
	}

	networkJSON, err := json.Marshal(configJSON.NetworkConfig)
	if err != nil {
		return nil, nil, err // un-tested
	}

	if configJSON.ProductProperties == nil {
		return nil, nil, fmt.Errorf("product properties must be provided in the config file")
	}

	propertiesJSON, err := json.Marshal(configJSON.ProductProperties)
	if err != nil {
		return nil, nil, err // un-tested
	}

	var minimalProperties *map[string]interface{}
	err = json.Unmarshal(propertiesJSON, &minimalProperties)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse product properties: %s", err)
	}

	combinedProperties := mergeProperties(*minimalProperties, additionalProperties)

	propertiesJSON, err = json.Marshal(combinedProperties)
	if err != nil {
		return nil, nil, err // un-tested
	}

	return propertiesJSON, networkJSON, nil
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
	configInput, err := ioutil.ReadFile(o.config.ConfigFile)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s (config file %s)", o.config.Name, err, o.config.ConfigFile)
	}

	propertiesJSON, networkJSON, err := extractPropertiesAndNetworkConfig(configInput, additionalProperties)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s (config file %s)", o.config.Name, err, o.config.ConfigFile)
	}

	err = o.omRunner.ResetAndConfigure(o.config.Name, o.config.Version, string(propertiesJSON), string(networkJSON))
	if err != nil {
		panic(err)
	}

	return nil
}
