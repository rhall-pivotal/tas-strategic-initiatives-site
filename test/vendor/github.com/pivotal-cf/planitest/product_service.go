package planitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

//go:generate counterfeiter -o ./fakes/command_runner.go --fake-name CommandRunner . CommandRunner
type CommandRunner interface {
	Run(string, ...string) (string, string, error)
}

type ProductConfig struct {
	Name         string
	Version      string
	ConfigFile   string
	MetadataFile string
}

type ConfigJSON struct {
	ProductProperties *map[string]interface{} `json:"product-properties,omitempty"`
	NetworkConfig     *map[string]interface{} `json:"network-config,omitempty"`
}

type ProductService struct {
	config        ProductConfig
	cmdRunner     CommandRunner
	RenderService RenderService
}

type StagedManifestResponse struct {
	Manifest map[string]interface{}
	Errors   OMError `json:"errors"`
}

type RenderService interface {
	RenderManifest(ProductConfig) (Manifest, error)
}

type OMError struct {
	// XXX: reconsider, the key here may change depending on the endpoint
	Messages []string `json:"base"`
}

func NewProductService(config ProductConfig) (*ProductService, error) {
	cmdRunner := NewExecutor()

	productService, err := NewProductServiceWithRunner(config, cmdRunner)
	if err != nil {
		return nil, err
	}

	return productService, err
}

func NewProductServiceWithRunner(config ProductConfig, cmdRunner CommandRunner) (*ProductService, error) {
	var renderService RenderService
	var err error

	switch os.Getenv("RENDERER") {
	case "om":
		renderService, err = NewOMServiceWithRunner(cmdRunner)
	case "ops-manifest":
		renderService, err = NewOpsManifestServiceWithRunner(cmdRunner)
	default:
		err = errors.New("RENDERER must be set to om or ops-manifest")
	}
	if err != nil {
		return nil, err
	}

	err = validateProductConfig(config)
	if err != nil {
		return nil, err
	}

	return &ProductService{config: config, cmdRunner: cmdRunner, RenderService: renderService}, nil
}

func (p *ProductService) Configure(additionalProperties map[string]interface{}) error {
	configInput, err := ioutil.ReadFile(p.config.ConfigFile)

	propertiesJSON, networkJSON, err := extractPropertiesAndNetworkConfig(configInput, additionalProperties)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s (config file %s)", p.config.Name, err, p.config.ConfigFile)
	}

	_, errOutput, err := p.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"revert-staged-changes",
	)

	if err != nil {
		return fmt.Errorf("Unable to revert staged changes: %s: %s", err, errOutput)
	}

	_, errOutput, err = p.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"stage-product",
		"--product-name", p.config.Name,
		"--product-version", p.config.Version,
	)

	if err != nil {
		return fmt.Errorf("Unable to stage product %q, version %q: %s: %s",
			p.config.Name, p.config.Version, err, errOutput)
	}

	_, errOutput, err = p.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"configure-product",
		"--product-name", p.config.Name,
		"--product-properties", string(propertiesJSON),
		"--product-network", string(networkJSON),
	)

	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s: %s", p.config.Name, err, errOutput)
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

func validateProductConfig(config ProductConfig) error {
	if len(config.Name) == 0 {
		return errors.New("Product name must be provided in config")
	}

	if len(config.Version) == 0 {
		return errors.New("Product version must be provided in config")
	}

	if len(config.ConfigFile) == 0 {
		return errors.New("Config file must be provided")
	}

	if os.Getenv("RENDERER") == "ops-manifest" {
		if len(config.MetadataFile) == 0 {
			return errors.New("Metadata file must be provided in config")
		}
	}

	return nil
}
