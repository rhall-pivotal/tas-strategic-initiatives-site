package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("loggregator agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "Pivotal Application Service for Windows"))
			Expect(tags).NotTo(HaveKey("product_version"))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		Context("when placement tags are configured by the user", func() {
			It("sets the placement tags on the emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".windows_diego_cell.placement_tags": "tag1,tag2",
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})
	})

	Describe("system metrics agent", func() {
		It("sets defaults on the system-metrics agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-system-metrics-agent")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := agent.Property("enabled")
			Expect(err).ToNot(HaveOccurred())
			Expect(enabled).To(BeTrue())

			tlsProps, err := agent.Property("system_metrics/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))
		})

		Context("when the Operator disables the system-metrics agent", func() {
			It("sets enabled to false", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.system_metrics_enabled": false,
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-system-metrics-agent")
				Expect(err).NotTo(HaveOccurred())

				enabled, err := agent.Property("enabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(enabled).To(BeFalse())

				tlsProps, err := agent.Property("system_metrics/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(tlsProps).To(HaveKey("ca_cert"))
				Expect(tlsProps).To(HaveKey("cert"))
				Expect(tlsProps).To(HaveKey("key"))
			})
		})
	})
})
