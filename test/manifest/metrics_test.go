package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metrics", func() {

	Describe("prom scraper", func() {
		It("configures the prom scraper", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			scraper, err := manifest.FindInstanceGroupJob("windows_diego_cell", "prom_scraper_windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(scraper)
		})
	})

	Describe("metrics discovery", func() {
		It("configures metrics-discovery on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			registrar, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metrics-discovery-registrar-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(registrar)

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metrics-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)
		})
	})
})
