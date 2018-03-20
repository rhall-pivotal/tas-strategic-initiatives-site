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
}

type ProductService struct {
	config    ProductConfig
	cmdRunner CommandRunner
}

func NewProductService(config ProductConfig) (*ProductService, error) {
	return NewProductServiceWithRunner(config, NewExecutor())
}

func NewProductServiceWithRunner(config ProductConfig, cmdRunner CommandRunner) (*ProductService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	err = validateProductConfig(config)
	if err != nil {
		return nil, err
	}

	return &ProductService{config: config, cmdRunner: cmdRunner}, nil
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

func (p *ProductService) RenderManifest() (Manifest, error) {
	response, errOutput, err := p.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"staged-manifest",
		"--product-name", p.config.Name,
	)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve staged manifest for product %q: %s: %s", p.config.Name, err, errOutput)
	}

	return NewManifest(string(response), p.cmdRunner), nil
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

func validateProductConfig(config ProductConfig) error {
	if len(config.Name) == 0 {
		return errors.New("Product name must be provided in config")
	}

	if len(config.Version) == 0 {
		return errors.New("Product version must be provided in config")
	}

	if len(config.PropertiesFile) == 0 {
		return errors.New("Properties file must be provided in config")
	}

	if len(config.NetworkConfigFile) == 0 {
		return errors.New("Network config file must be provided in config")
	}

	return nil
}
