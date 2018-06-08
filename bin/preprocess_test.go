package main_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("preprocess-metadata-parts", func() {
	var (
		outputPath        string
		metadataPartsPath string
	)

	BeforeEach(func() {
		var err error
		outputPath, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		metadataPartsPath = filepath.Join("fixtures", "valid")
	})

	It("processes the templates files for the ERT", func() {
		command := exec.Command(pathToMain,
			"--tile-name", "ert",
			"--input-path", metadataPartsPath,
			"--output-path", outputPath,
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		baseFilePath := filepath.Join(outputPath, "base.yml")
		contents, err := ioutil.ReadFile(baseFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
metadata_version: some-metadata-version
name: ert
provides_product_versions:
- name: ert-product
requires_product_versions:
- name: some-other-product
  version: 1.2.3.4
product_version: some-product-version
minimum_version_for_upgrade: some-minimum-version
label: some-label
description: some-description
icon_image: some-icon
rank: 90
serial: false
post_deploy_errands:
  - name: some-errand
variables:
- name: root-ca
  type: rsa
  options:
    is_ca: true
`))

		ertTemplateFilePath := filepath.Join(outputPath, "ert_jobs", "template.yml")
		contents, err = ioutil.ReadFile(ertTemplateFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
name: ert-job-template
label: Job for ERT
count: 3
configurable: true
description:
- Multi-line content
- for the ERT
`))

		srtTemplateFilePath := filepath.Join(outputPath, "srt_jobs", "template.yml")
		contents, err = ioutil.ReadFile(srtTemplateFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
name: srt-job-template
label: Job for ERT
count: 0
configurable: false
`))
	})

	It("processes the templates files for the SRT", func() {
		command := exec.Command(pathToMain,
			"--tile-name", "srt",
			"--input-path", metadataPartsPath,
			"--output-path", outputPath,
		)

		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
		Eventually(session).Should(gexec.Exit(0))

		baseFilePath := filepath.Join(outputPath, "base.yml")
		contents, err := ioutil.ReadFile(baseFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
metadata_version: some-metadata-version
name: srt
provides_product_versions:
- name: srt-product
requires_product_versions:
- name: some-other-product
  version: 1.2.3.4
product_version: some-product-version
minimum_version_for_upgrade: some-minimum-version
label: some-label
description: some-description
icon_image: some-icon
rank: 90
serial: false
post_deploy_errands:
  - name: some-errand
variables:
- name: root-ca
  type: rsa
  options:
    is_ca: true
`))

		ertTemplateFilePath := filepath.Join(outputPath, "ert_jobs", "template.yml")
		contents, err = ioutil.ReadFile(ertTemplateFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
name: ert-job-template
label: Job for SRT
count: 0
configurable: false
`))

		srtTemplateFilePath := filepath.Join(outputPath, "srt_jobs", "template.yml")
		contents, err = ioutil.ReadFile(srtTemplateFilePath)
		Expect(err).NotTo(HaveOccurred())
		Expect(contents).To(MatchYAML(`---
name: srt-job-template
label: Job for SRT
count: 1
configurable: true
description:
- Multi-line content
- for the SRT
`))
	})

	Context("failure cases", func() {
		Context("when the metadata file contains a bad template definition", func() {
			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--tile-name", "ert",
					"--input-path", filepath.Join("fixtures", "bad-template"),
					"--output-path", outputPath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("unclosed action"))
			})
		})

		Context("when the --tile-name flag is not provided", func() {
			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--input-path", metadataPartsPath,
					"--output-path", outputPath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("please provide a tile name using the --tile-name option"))
			})
		})

		Context("when the --input-path flag is not provided", func() {
			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--tile-name", "ert",
					"--output-path", outputPath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("please provide a metadata parts directory path using the --input-path option"))
			})
		})

		Context("when the --output-path flag is not provided", func() {
			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--tile-name", "ert",
					"--input-path", metadataPartsPath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("please provide an output directory path using the --output-path option"))
			})
		})

		Context("when the output directory path is actually a file", func() {
			var existingFilePath string

			BeforeEach(func() {
				existingFile, err := ioutil.TempFile("", "")
				Expect(err).NotTo(HaveOccurred())

				existingFilePath = existingFile.Name()
			})

			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--tile-name", "ert",
					"--input-path", metadataPartsPath,
					"--output-path", existingFilePath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("not a directory"))
			})
		})

		Context("when an unsupported tile name is specified", func() {
			It("prints an error message", func() {
				command := exec.Command(pathToMain,
					"--tile-name", "some-other-tile",
					"--input-path", metadataPartsPath,
					"--output-path", outputPath,
				)

				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(session.Err.Contents()).To(ContainSubstring("unsupported tile name: some-other-tile"))
			})
		})
	})
})
