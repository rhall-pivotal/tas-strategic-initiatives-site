package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Logging", func() {
	Describe("loggregator agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			v2Api, err := agent.Property("loggregator/use_v2_api")
			Expect(err).ToNot(HaveOccurred())
			Expect(v2Api).To(BeTrue())

			tlsProps, err := agent.Property("loggregator/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))

			expectSecureMetrics(agent)

			tlsAgentProps, err := agent.Property("loggregator/tls/agent")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsAgentProps).To(HaveKey("cert"))
			Expect(tlsAgentProps).To(HaveKey("key"))

			By("disabling udp")
			udpDisabled, err := agent.Property("disable_udp")
			Expect(err).NotTo(HaveOccurred())
			Expect(udpDisabled).To(BeTrue())

			By("getting the grpc port")
			port, err := agent.Property("grpc_port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3459))

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "Pivotal Application Service for Windows"))
			Expect(tags).NotTo(HaveKey("product_version"))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		It("is enabled by default", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			_, err = agent.Property("loggregator_agent/enabled")
			Expect(err).ToNot(HaveOccurred())
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

	Describe("forwarder agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-forwarder-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)

			By("getting the grpc port")
			port, err := agent.Property("port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3458))

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

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-forwarder-agent-windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})
	})

	Describe("syslog agent", func() {
		It("sets defaults on the syslog agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggr-syslog-agent-windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(agent)

			port, err := agent.Property("port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(3460))

			tlsProps, err := agent.Property("tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))

			cacheTlsProps, err := agent.Property("cache/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(cacheTlsProps).To(HaveKey("ca_cert"))
			Expect(cacheTlsProps).To(HaveKey("cert"))
			Expect(cacheTlsProps).To(HaveKey("key"))
			Expect(cacheTlsProps).To(HaveKeyWithValue("cn", "binding-cache"))
		})
	})

	Describe("prom scraper", func() {
		It("configures the prom scraper", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			scraper, err := manifest.FindInstanceGroupJob("windows_diego_cell", "prom_scraper_windows")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(scraper)
		})
	})
})

func expectSecureMetrics(job planitest.Manifest) {
	metricsProps, err := job.Property("metrics")
	Expect(err).ToNot(HaveOccurred())
	Expect(metricsProps).To(HaveKey("ca_cert"))
	Expect(metricsProps).To(HaveKey("cert"))
	Expect(metricsProps).To(HaveKey("key"))
	Expect(metricsProps).To(HaveKey("server_name"))
}