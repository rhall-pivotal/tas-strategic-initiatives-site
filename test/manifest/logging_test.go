package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {

	var instanceGroups []string = []string{"isolated_diego_cell", "isolated_ha_proxy", "isolated_router"}

	Describe("loggregator agent", func() {
		It("sets tags on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(tags).To(HaveKeyWithValue("placement_tag", "isosegtag"))
				Expect(tags).To(HaveKeyWithValue("product", "PCF Isolation Segment"))
				Expect(tags).NotTo(HaveKey("product_version"))
				Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
			}
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
			}
		})
	})
})
