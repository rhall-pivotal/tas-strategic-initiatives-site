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
	Name              string
	Version           string
	PropertiesFile    string
	NetworkConfigFile string
	ConfigFile        string
	MetadataFile      string
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
	err := validateProductConfig(config)
	if err != nil {
		return nil, err
	}

	var renderService RenderService
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

	return &ProductService{config: config, cmdRunner: cmdRunner, RenderService: renderService}, nil
}

func (p *ProductService) Configure(additionalProperties map[string]interface{}) error {

	propertiesJSON, err := ioutil.ReadFile(p.config.PropertiesFile)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s", p.config.Name, err)
	}

	var minimalProperties map[string]interface{}
	err = json.Unmarshal(propertiesJSON, &minimalProperties)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: could not parse properties file %q: %s", p.config.Name, p.config.PropertiesFile, err)
	}

	networkJSON, err := ioutil.ReadFile(p.config.NetworkConfigFile)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s", p.config.Name, err)
	}

	combinedProperties := mergeProperties(minimalProperties, additionalProperties)

	propertiesJSON, err = json.Marshal(combinedProperties)
	if err != nil {
		return fmt.Errorf("Unable to configure product %q: %s", p.config.Name, err) // un-tested
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

	if os.Getenv("RENDERER") == "ops-manifest" {
		if len(config.MetadataFile) == 0 {
			return errors.New("Metadata file must be provided in config")
		}
		if len(config.ConfigFile) == 0 {
			return errors.New("Config file must be provided")
		}
	} else {
		if len(config.PropertiesFile) == 0 {
			return errors.New("Properties file must be provided in config")
		}

		if len(config.NetworkConfigFile) == 0 {
			return errors.New("Network config file must be provided in config")
		}

	}

	return nil
}
