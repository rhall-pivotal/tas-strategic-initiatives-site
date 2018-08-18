package manifest_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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
var configFile *os.File

var _ = SynchronizedBeforeSuite(func() []byte {
	cmd := exec.Command("../../bin/build")
	cmd.Env = append(os.Environ(),
		"METADATA_ONLY=true",
		"STUB_RELEASES=true",
	)

	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running bin/build: %s", err.Error())
	}
	Expect(err).NotTo(HaveOccurred())

	return []byte(output)
}, func(metadataContents []byte) {
	metadataFile = string(metadataContents)
	var err error

	configFile, err = os.Open("config.json")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	configFile.Close()
})

var _ = BeforeEach(func() {
	var err error
	productConfig = planitest.ProductConfig{
		ConfigFile: configFile,
		TileFile:   strings.NewReader(metadataFile),
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
