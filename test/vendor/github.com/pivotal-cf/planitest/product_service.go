package planitest

import (
	"errors"
	"os"

	"github.com/pivotal-cf/planitest/internal"
)

type ProductConfig struct {
	Name         string
	Version      string
	ConfigFile   string
	MetadataFile string
}

type ConfigJSON struct {
	ProductProperties *map[string]interface{} `json:"product-properties,omitempty"`
	NetworkProperties *map[string]interface{} `json:"network-properties,omitempty"`
}

type ProductService struct {
	config        ProductConfig
	RenderService RenderService
}

type StagedManifestResponse struct {
	Manifest map[string]interface{}
	Errors   OMError `json:"errors"`
}

type RenderService interface {
	RenderManifest(additionalProperties map[string]interface{}) (Manifest, error)
}

type OMError struct {
	// XXX: reconsider, the key here may change depending on the endpoint
	Messages []string `json:"base"`
}

func NewProductService(config ProductConfig) (*ProductService, error) {
	omRunner := internal.NewOMRunner(NewExecutor())
	opsManifestRunner := internal.NewOpsManifestRunner(NewExecutor())

	productService, err := NewProductServiceWithRunner(config, omRunner, opsManifestRunner)
	if err != nil {
		return nil, err
	}

	return productService, err
}

func NewProductServiceWithRunner(config ProductConfig, omRunner OMRunner, opsManifestRunner OpsManifestRunner) (*ProductService, error) {
	var renderService RenderService
	var err error

	switch os.Getenv("RENDERER") {
	case "om":
		renderService, err = NewOMServiceWithRunner(OMConfig{
			Name:       config.Name,
			Version:    config.Version,
			ConfigFile: config.ConfigFile,
		}, omRunner)
	case "ops-manifest":
		renderService, err = NewOpsManifestServiceWithRunner(OpsManifestConfig{
			ConfigFile:   config.ConfigFile,
			MetadataFile: config.MetadataFile,
		}, opsManifestRunner)
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

	return &ProductService{config: config, RenderService: renderService}, nil
}

func validateProductConfig(config ProductConfig) error {
	if len(config.ConfigFile) == 0 {
		return errors.New("Config file must be provided")
	}

	if os.Getenv("RENDERER") == "om" {
		if len(config.Name) == 0 {
			return errors.New("Product name must be provided in config")
		}

		if len(config.Version) == 0 {
			return errors.New("Product version must be provided in config")
		}
	}

	if os.Getenv("RENDERER") == "ops-manifest" {
		if len(config.MetadataFile) == 0 {
			return errors.New("Metadata file must be provided in config")
		}
	}

	return nil
}
