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

			d, err := loadDomain("../../properties/metrics.yml", "prom_scraper_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err := scraper.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))
		})
	})

	Describe("metrics discovery", func() {
		It("configures metrics-discovery on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			registrar, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metrics-discovery-registrar-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(registrar)

			d, err := loadDomain("../../properties/metrics.yml", "metrics_discovery_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err := registrar.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metrics-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)

			d, err = loadDomain("../../properties/metrics.yml", "metrics_agent_metrics_tls")
			Expect(err).ToNot(HaveOccurred())

			metricsProps, err = agent.Property("metrics")
			Expect(err).ToNot(HaveOccurred())
			Expect(metricsProps).To(HaveKeyWithValue("server_name", d))
		})
	})
})
