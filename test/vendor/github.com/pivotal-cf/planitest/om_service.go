package planitest

import (
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

func (o OMService) RenderManifest(additionalProperties map[string]interface{}) (Manifest, error) {
	stagedProduct, err := o.omRunner.FindStagedProduct(o.config.Name)
	if err != nil {
		return Manifest{}, err
	}

	configInput, err := MergeAdditionalProductProperties(o.config.ConfigFile, additionalProperties)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to configure product %q: %s", o.config.Name, err)
	}

	err = o.omRunner.ResetAndConfigure(o.config.Name, o.config.Version, string(configInput))
	if err != nil {
		return Manifest{}, err
	}

	// calling configure can re-staged and update the product GUID
	stagedProduct, err = o.omRunner.FindStagedProduct(o.config.Name)
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
