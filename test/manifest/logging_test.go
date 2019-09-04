package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Logging", func() {
	var (
		instanceGroups       []string = []string{"isolated_diego_cell", "isolated_ha_proxy", "isolated_router"}
		getAllInstanceGroups func(planitest.Manifest) []string
	)

	getAllInstanceGroups = func(manifest planitest.Manifest) []string {
		groups, err := manifest.Path("/instance_groups")
		Expect(err).NotTo(HaveOccurred())

		groupList, ok := groups.([]interface{})
		Expect(ok).To(BeTrue())

		names := []string{}
		for _, group := range groupList {
			groupName := group.(map[interface{}]interface{})["name"].(string)

			// ignore VMs that only contain a single placeholder job, i.e. SF-PAS only VMs that are present but non-configurable in PAS build
			jobs, err := manifest.Path(fmt.Sprintf("/instance_groups/name=%s/jobs", groupName))
			Expect(err).NotTo(HaveOccurred())
			if len(jobs.([]interface{})) > 1 {
				names = append(names, groupName)
			}
		}
		Expect(names).NotTo(BeEmpty())
		return names
	}

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

				cacheTlsProps, err := agent.Property("cache/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(cacheTlsProps).To(HaveKey("ca_cert"))
				Expect(cacheTlsProps).To(HaveKey("cert"))
				Expect(cacheTlsProps).To(HaveKey("key"))
				Expect(cacheTlsProps).To(HaveKeyWithValue("cn", "binding-cache"))
			}
		})
	})

	Describe("prom scraper", func() {
		It("configures the prom scraper on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				_, err := manifest.FindInstanceGroupJob(ig, "prom_scraper")
				Expect(err).NotTo(HaveOccurred())
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
					".properties.system_logging":              "enabled",
					".properties.system_logging.enabled.host": "example.com",
					".properties.system_logging.enabled.port": 2514,
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

	Describe("system metrics agent", func() {
		It("sets defaults on the system-metrics agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggr-system-metrics-agent")
				Expect(err).NotTo(HaveOccurred())

				enabled, err := agent.Property("enabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(enabled).To(BeTrue())

				tlsProps, err := agent.Property("system_metrics/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(tlsProps).To(HaveKey("ca_cert"))
				Expect(tlsProps).To(HaveKey("cert"))
				Expect(tlsProps).To(HaveKey("key"))
			}
		})

		Context("when the Operator disables the system-metrics agent", func() {
			It("sets enabled to false", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.system_metrics_enabled": false,
				})
				Expect(err).NotTo(HaveOccurred())

				instanceGroups := getAllInstanceGroups(manifest)

				for _, ig := range instanceGroups {
					agent, err := manifest.FindInstanceGroupJob(ig, "loggr-system-metrics-agent")
					Expect(err).NotTo(HaveOccurred())

					enabled, err := agent.Property("enabled")
					Expect(err).ToNot(HaveOccurred())
					Expect(enabled).To(BeFalse())

					tlsProps, err := agent.Property("system_metrics/tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(tlsProps).To(HaveKey("ca_cert"))
					Expect(tlsProps).To(HaveKey("cert"))
					Expect(tlsProps).To(HaveKey("key"))
				}
			})
		})
	})
})
