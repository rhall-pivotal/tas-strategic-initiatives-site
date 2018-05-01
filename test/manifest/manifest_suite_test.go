package manifest_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"

	"testing"
)

func TestManifestGeneration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Manifest Generation Suite")
}

var product *planitest.ProductService
var productConfig planitest.ProductConfig

var _ = BeforeEach(func() {
	productVersion, err := fetchProductVersion()
	Expect(err).NotTo(HaveOccurred())

	productConfig = planitest.ProductConfig{
		Name:              "cf",
		Version:           productVersion,
		ConfigFile:        "product_properties.yml",
		PropertiesFile:    "product_config.json",
		MetadataFile:      "cf-metadata.yml",
		NetworkConfigFile: "network_config.json",
	}
	product, err = planitest.NewProductService(productConfig)
	Expect(err).NotTo(HaveOccurred())
})

func fetchProductVersion() (string, error) {
	contents, err := ioutil.ReadFile(filepath.Join("..", "..", "version"))
	if err != nil {
		return "", err
	}

	matches := regexp.MustCompile(`(\d\.\d{1,2}\.\d{1,3})\-build\.\d{1,3}`).FindStringSubmatch(string(contents))

	if len(matches) != 2 {
		return "", fmt.Errorf("could not find version in %q", contents)
	}

	return matches[1], nil
}
