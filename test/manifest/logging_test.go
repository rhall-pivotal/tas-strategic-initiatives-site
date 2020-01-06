package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Logging", func() {
	var (
		instanceGroups []string = []string{"isolated_diego_cell", "isolated_ha_proxy", "isolated_router"}
	)

	Describe("loggregator agent", func() {
		It("sets tags on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
				Expect(err).NotTo(HaveOccurred())

				tlsProps, err := agent.Property("loggregator/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(tlsProps).To(HaveKey("ca_cert"))

				expectSecureMetrics(agent)

				tlsAgentProps, err := agent.Property("loggregator/tls/agent")
				Expect(err).ToNot(HaveOccurred())
				Expect(tlsAgentProps).To(HaveKey("cert"))
				Expect(tlsAgentProps).To(HaveKey("key"))

				port, err := agent.Property("grpc_port")
				Expect(err).NotTo(HaveOccurred())
				Expect(port).To(Equal(3459))

				udpDisabled, err := agent.Property("disable_udp")
				Expect(err).NotTo(HaveOccurred())
				Expect(udpDisabled).To(BeTrue())
			}
		})

		It("is enabled by default", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
				Expect(err).NotTo(HaveOccurred())

				_, err = agent.Property("loggregator_agent/enabled")
				Expect(err).ToNot(HaveOccurred())
			}
		})

		Describe("tags", func() {
			Context("when compute isolation is enabled", func() {
				It("adds the appropriate manifest for tags", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					for _, ig := range instanceGroups {
						agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
						Expect(err).NotTo(HaveOccurred())

						tags, err := agent.Property("tags")
						Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
						Expect(tags).To(HaveKeyWithValue("placement_tag", "isosegtag"))
						Expect(tags).To(HaveKeyWithValue("product", "Pivotal Isolation Segment"))
						Expect(tags).NotTo(HaveKey("product_version"))
						Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
					}
				})
			})

			Context("when compute isolation is disabled", func() {
				It("adds the appropriate manifest for tags", func() {
					manifest, err := product.RenderManifest(map[string]interface{}{
						".properties.compute_isolation":                                "disabled",
						".properties.compute_isolation.enabled.isolation_segment_name": "",
					})
					Expect(err).NotTo(HaveOccurred())

					for _, ig := range instanceGroups {
						agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
						Expect(err).NotTo(HaveOccurred())

						tags, err := agent.Property("tags")
						Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
						Expect(tags).NotTo(HaveKey("placement_tag"))
						Expect(tags).To(HaveKeyWithValue("product", "Pivotal Isolation Segment"))
						Expect(tags).NotTo(HaveKey("product_version"))
						Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
					}
				})
			})
		})
	})

	Describe("syslog agent", func() {
		It("sets defaults on the syslog agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggr-syslog-agent")
				Expect(err).NotTo(HaveOccurred())

				port, err := agent.Property("port")
				Expect(err).NotTo(HaveOccurred())
				Expect(port).To(Equal(3460))

				tlsProps, err := agent.Property("tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(tlsProps).To(HaveKey("ca_cert"))
				Expect(tlsProps).To(HaveKey("cert"))
				Expect(tlsProps).To(HaveKey("key"))

				expectSecureMetrics(agent)

				cacheTlsProps, err := agent.Property("cache/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(cacheTlsProps).To(HaveKey("ca_cert"))
				Expect(cacheTlsProps).To(HaveKey("cert"))
				Expect(cacheTlsProps).To(HaveKey("key"))
				Expect(cacheTlsProps).To(HaveKeyWithValue("cn", "binding-cache"))
			}
		})

		It("has aggreate drain url", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggr-syslog-agent")
				Expect(err).NotTo(HaveOccurred())

				aggregateDrains, err := agent.Property("aggregate_drains")
				Expect(err).NotTo(HaveOccurred())
				Expect(aggregateDrains).To(ContainSubstring("syslog-tls://doppler.service.cf.internal:6067"))
			}
		})
	})

	Describe("prom scraper", func() {
		It("configures the prom scraper on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				scraper, err := manifest.FindInstanceGroupJob(ig, "prom_scraper")
				Expect(err).NotTo(HaveOccurred())

				expectSecureMetrics(scraper)
			}
		})

		Describe("forwarder agent", func() {
			It("sets defaults on the forwarder agent", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				for _, ig := range instanceGroups {
					agent, err := manifest.FindInstanceGroupJob(ig, "loggr-forwarder-agent")
					Expect(err).NotTo(HaveOccurred())

					By("getting the grpc port")
					port, err := agent.Property("port")
					Expect(err).NotTo(HaveOccurred())
					Expect(port).To(Equal(3458))

					expectSecureMetrics(agent)
				}
			})

			Describe("tags", func() {
				Context("when compute isolation is enabled", func() {
					It("adds the appropriate manifest for tags", func() {
						manifest, err := product.RenderManifest(nil)
						Expect(err).NotTo(HaveOccurred())

						for _, ig := range instanceGroups {
							agent, err := manifest.FindInstanceGroupJob(ig, "loggr-forwarder-agent")
							Expect(err).NotTo(HaveOccurred())

							tags, err := agent.Property("tags")
							Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
							Expect(tags).To(HaveKeyWithValue("placement_tag", "isosegtag"))
							Expect(tags).To(HaveKeyWithValue("product", "Pivotal Isolation Segment"))
							Expect(tags).NotTo(HaveKey("product_version"))
							Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
						}
					})
				})

				Context("when compute isolation is disabled", func() {
					It("adds the appropriate manifest for tags", func() {
						manifest, err := product.RenderManifest(map[string]interface{}{
							".properties.compute_isolation":                                "disabled",
							".properties.compute_isolation.enabled.isolation_segment_name": "",
						})
						Expect(err).NotTo(HaveOccurred())

						for _, ig := range instanceGroups {
							agent, err := manifest.FindInstanceGroupJob(ig, "loggr-forwarder-agent")
							Expect(err).NotTo(HaveOccurred())

							tags, err := agent.Property("tags")
							Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
							Expect(tags).NotTo(HaveKey("placement_tag"))
							Expect(tags).To(HaveKeyWithValue("product", "Pivotal Isolation Segment"))
							Expect(tags).NotTo(HaveKey("product_version"))
							Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
						}
					})
				})
			})
		})
	})

	Describe("syslog forwarding", func() {

		It("includes the vcap rule", func() {
			for _, ig := range instanceGroups {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.system_logging":                  "enabled",
					".properties.system_logging.enabled.host":     "example.com",
					".properties.system_logging.enabled.port":     2514,
					".properties.system_logging.enabled.protocol": "tcp",
				})
				Expect(err).NotTo(HaveOccurred())

				syslogForwarder, err := manifest.FindInstanceGroupJob(ig, "syslog_forwarder")
				Expect(err).NotTo(HaveOccurred())

				syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
				Expect(err).NotTo(HaveOccurred())
				Expect(syslogConfig).To(ContainSubstring(`if ($programname startswith "vcap.") then stop`))
			}
		})

		Context("when a custom rule is specified", func() {
			It("adds the custom rule", func() {
				multilineRule := `
some
multi
line
rule
`
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.system_logging":                     "enabled",
					".properties.system_logging.enabled.host":        "example.com",
					".properties.system_logging.enabled.port":        2514,
					".properties.system_logging.enabled.syslog_rule": multilineRule,
					".properties.system_logging.enabled.protocol":    "tcp",
				})
				Expect(err).NotTo(HaveOccurred())

				for _, ig := range instanceGroups {
					syslogForwarder, err := manifest.FindInstanceGroupJob(ig, "syslog_forwarder")
					Expect(err).NotTo(HaveOccurred())

					syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
					Expect(err).NotTo(HaveOccurred())
					Expect(syslogConfig).To(ContainSubstring(`
some
multi
line
rule
`))
				}
			})
		})

		Context("when dropping debug logs", func() {
			It("does not forward debug logs", func() {
				for _, ig := range instanceGroups {
					manifest, err := product.RenderManifest(map[string]interface{}{
						".properties.system_logging":                           "enabled",
						".properties.system_logging.enabled.host":              "example.com",
						".properties.system_logging.enabled.port":              2514,
						".properties.system_logging.enabled.syslog_drop_debug": true,
						".properties.system_logging.enabled.protocol":          "tcp",
					})
					Expect(err).NotTo(HaveOccurred())

					syslogForwarder, err := manifest.FindInstanceGroupJob(ig, "syslog_forwarder")
					Expect(err).NotTo(HaveOccurred())

					syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
					Expect(err).NotTo(HaveOccurred())

					Expect(syslogConfig).To(ContainSubstring(`if ($msg contains "DEBUG") then stop`))
				}
			})
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
