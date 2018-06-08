package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestPreprocessMetadata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "preprocess-metadata")
}

var pathToMain string

var _ = BeforeSuite(func() {
	var err error
	pathToMain, err = gexec.Build("github.com/pivotal-cf/p-runtime/bin")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
