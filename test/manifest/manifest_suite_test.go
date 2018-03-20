package manifest_test

import (
	"fmt"
	"io/ioutil"
	"os"
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

var _ = BeforeEach(func() {
	productVersion, err := fetchProductVersion()
	Expect(err).NotTo(HaveOccurred())

	product, err = planitest.NewProductService(planitest.ProductConfig{
		Name:              "cf",
		Version:           productVersion,
		PropertiesFile:    "product_config.json",
		NetworkConfigFile: os.Getenv("NETWORK_CONFIG_FILE"),
	})
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
