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
var productConfig planitest.ProductConfig
var metadataFile string

var _ = SynchronizedBeforeSuite(func() []byte {
	var env []string

	file, err := ioutil.TempFile("", "Planitest")
	Expect(err).NotTo(HaveOccurred())

	env = append(env, "METADATA_ONLY=true", "STUB_RELEASES=true")

	cmd := planitest.NewExecutorWithEnv(env)

	output, errOutput, err := cmd.Run("../../bin/build")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running bin/build: %s: %s", err, errOutput)
	}
	Expect(err).NotTo(HaveOccurred())

	_, err = file.WriteString(output)
	Expect(err).NotTo(HaveOccurred())

	return []byte(file.Name())
}, func(path []byte) {
	metadataFile = string(path)
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	err := os.Remove(metadataFile)
	Expect(err).NotTo(HaveOccurred())
})

var _ = BeforeEach(func() {
	productVersion, err := fetchProductVersion()
	Expect(err).NotTo(HaveOccurred())

	productConfig = planitest.ProductConfig{
		Name:         "p-isolation-segment",
		Version:      productVersion,
		ConfigFile:   "config.json",
		MetadataFile: metadataFile,
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
