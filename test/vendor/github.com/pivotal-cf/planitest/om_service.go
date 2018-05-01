package planitest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type StagedProduct struct {
	GUID           string `json:"guid"`
	Type           string `json:"type"`
	ProductVersion string `json:"product_version"`
}

type OMService struct {
	cmdRunner CommandRunner
}

func NewOMService() (*OMService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	return NewOMServiceWithRunner(NewExecutor())
}

func NewOMServiceWithRunner(cmdRunner CommandRunner) (*OMService, error) {
	err := validateEnvironmentVariables()
	if err != nil {
		return nil, err
	}

	return &OMService{cmdRunner: cmdRunner}, nil
}

func (o OMService) StagedProducts() ([]StagedProduct, error) {
	response, errOutput, err := o.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"curl",
		"--path", "/api/v0/staged/products",
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve staged products: %s: %s", err, errOutput)
	}

	var stagedProducts []StagedProduct
	err = json.Unmarshal([]byte(response), &stagedProducts)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve staged products: %s", err)
	}

	return stagedProducts, nil
}

func (o OMService) FindStagedProduct(productName string) (StagedProduct, error) {
	stagedProducts, _ := o.StagedProducts()

	var stagedTypes []string
	for _, sp := range stagedProducts {
		if sp.Type == productName {
			return sp, nil
		} else {
			stagedTypes = append(stagedTypes, sp.Type)
		}
	}

	return StagedProduct{}, fmt.Errorf("Product %q has not been staged. Staged products: %q",
		productName, strings.Join(stagedTypes, ", "))
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

func (o OMService) RenderManifest(config ProductConfig) (Manifest, error) {
	stagedProduct, err := o.FindStagedProduct(config.Name)
	if err != nil {
		return Manifest{}, err
	}

	response, errOutput, err := o.cmdRunner.Run(
		"om",
		"--skip-ssl-validation",
		"--target", os.Getenv("OM_URL"),
		"curl",
		"--path", fmt.Sprintf("/api/v0/staged/products/%s/manifest", stagedProduct.GUID),
	)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve staged manifest for product guid %q: %s: %s", stagedProduct.GUID, err, errOutput)
	}
	var smr StagedManifestResponse
	err = json.Unmarshal([]byte(response), &smr)
	if err != nil {
		return Manifest{}, fmt.Errorf("Unable to retrieve staged manifest for product guid %q: %s", stagedProduct.GUID, err)
	}
	if len(smr.Errors.Messages) > 0 {
		return Manifest{}, fmt.Errorf("Unable to retrieve staged manifest for product guid %q: %s",
			stagedProduct.GUID,
			smr.Errors.Messages[0])
	}

	y, err := yaml.Marshal(smr.Manifest)
	if err != nil {
		return Manifest{}, err // un-tested
	}

	return NewManifest(string(y), o.cmdRunner), nil
}
