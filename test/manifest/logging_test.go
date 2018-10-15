package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("loggregator agent", func() {
		var (
			instanceGroups []string
			productTag     string
		)

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroups = []string{"blobstore", "compute", "control", "database"}
				productTag = "Small Footprint PAS"
			} else {
				instanceGroups = []string{
					"clock_global",
					"cloud_controller",
					"cloud_controller_worker",
					"diego_brain",
					"diego_cell",
					"diego_database",
					"doppler",
					"ha_proxy",
					"loggregator_trafficcontroller",
					"mysql",
					"nats",
					"nfs_server",
					"router",
					"syslog_adapter",
					"syslog_scheduler",
					"tcp_router",
					"uaa",
				}
				productTag = "Pivotal Application Service"
			}
		})

		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
				Expect(err).NotTo(HaveOccurred())

				By("disabling the cf deployment name in emitted metrics")
				deploymentName, err := agent.Property("deployment")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(deploymentName).To(Equal(""), "Instance Group: %s", ig)

				By("adding tags to the metrics emitted")
				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(tags).To(HaveKeyWithValue("product", productTag))
				Expect(tags).To(HaveKeyWithValue("product_version", MatchRegexp(`^\d+\.\d+\.\d+.*`)))
				Expect(tags).To(HaveKeyWithValue("system_domain", "sys.example.com"))
			}
		})

		Context("when the enable cf metric name is set to true (migration during upgrades)", func() {
			It("sets the metric deployment name to cf", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.enable_cf_metric_name": true,
				})
				Expect(err).NotTo(HaveOccurred())

				for _, ig := range instanceGroups {
					agent, err := manifest.FindInstanceGroupJob(ig, "loggregator_agent")
					Expect(err).NotTo(HaveOccurred())

					deploymentName, err := agent.Property("deployment")
					Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
					Expect(deploymentName).To(Equal("cf"), "Instance Group: %s", ig)
				}
			})
		})
	})

	Describe("log cache", func() {
		var instanceGroup string
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "doppler"
			}
		})

		It("is enabled by default", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
			Expect(err).NotTo(HaveOccurred())

			disabledProperty, err := logCache.Property("disabled")
			Expect(err).ToNot(HaveOccurred())
			Expect(disabledProperty).To(BeFalse())
		})

		It("has tls server certs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
			Expect(err).NotTo(HaveOccurred())

			tlsProps, err := logCache.Property("tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))
		})

		It("specifies the port to listen on", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
			Expect(err).NotTo(HaveOccurred())

			port, err := logCache.Property("port")
			Expect(err).ToNot(HaveOccurred())

			if productName == "srt" {
				Expect(port).To(Equal(8090))
			} else {
				Expect(port).To(Equal(8080))
			}
		})

		It("has a log-cache-gateway with a gateway addr", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-gateway")
			Expect(err).NotTo(HaveOccurred())

			gatewayAddr, err := logCache.Property("gateway_addr")
			Expect(err).ToNot(HaveOccurred())
			if productName == "srt" {
				Expect(gatewayAddr).To(Equal("localhost:8087"))
			} else {
				Expect(gatewayAddr).To(Equal("localhost:8081"))
			}
		})

		It("has a log-cache-nozzle with tls certs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			nozzle, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-nozzle")
			Expect(err).NotTo(HaveOccurred())

			tlsProps, err := nozzle.Property("logs_provider/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))
		})

		It("has a log-cache-expvar-forwarder job with templated counters/gauges", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			forwarder, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-expvar-forwarder")
			Expect(err).NotTo(HaveOccurred())

			counters, err := forwarder.Property("counters")
			Expect(err).ToNot(HaveOccurred())
			Expect(counters).To(ContainElement(map[interface{}]interface{}{
				"addr":      "http://localhost:6060/debug/vars",
				"name":      "egress",
				"source_id": "log-cache",
				"template":  "{{.LogCache.Egress}}",
			}))

			gauges, err := forwarder.Property("gauges")
			Expect(err).ToNot(HaveOccurred())
			Expect(gauges).To(ContainElement(map[interface{}]interface{}{
				"addr":      "http://localhost:6060/debug/vars",
				"name":      "cache-period",
				"source_id": "log-cache",
				"template":  "{{.LogCache.CachePeriod}}",
			}))
		})

		It("registers the log-cache route", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			routeRegistrar, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
			Expect(err).NotTo(HaveOccurred())

			routes, err := routeRegistrar.Property("route_registrar/routes")
			Expect(err).ToNot(HaveOccurred())
			Expect(routes).To(ContainElement(HaveKeyWithValue("uris", []interface{}{
				"log-cache.sys.example.com",
			})))

			if productName == "srt" {
				Expect(routes).To(ContainElement(HaveKeyWithValue("port", 8089)))
			} else {
				Expect(routes).To(ContainElement(HaveKeyWithValue("port", 8083)))
			}
		})

		It("has an auth proxy", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			proxy, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-cf-auth-proxy")
			Expect(err).NotTo(HaveOccurred())

			ccProperties, err := proxy.Property("cc")
			Expect(err).ToNot(HaveOccurred())

			Expect(ccProperties).To(HaveKeyWithValue(
				"common_name", "cloud-controller-ng.service.cf.internal"))
			Expect(ccProperties).To(HaveKeyWithValue(
				"capi_internal_addr", "https://cloud-controller-ng.service.cf.internal:9023"))

			Expect(ccProperties).To(HaveKey("ca_cert"))
			Expect(ccProperties).To(HaveKey("cert"))
			Expect(ccProperties).To(HaveKey("key"))

			proxyPort, err := proxy.Property("proxy_port")
			Expect(err).ToNot(HaveOccurred())

			if productName == "srt" {
				Expect(proxyPort).To(Equal(8089))
			} else {
				Expect(proxyPort).To(Equal(8083))
			}

			uaaProperties, err := proxy.Property("uaa")
			Expect(err).ToNot(HaveOccurred())

			Expect(uaaProperties).To(HaveKeyWithValue("client_id", "doppler"))
			Expect(uaaProperties).To(HaveKeyWithValue("internal_addr", "https://uaa.service.cf.internal:8443"))

			Expect(uaaProperties).To(HaveKey("ca_cert"))
			Expect(uaaProperties).To(HaveKey("client_secret"))

		})

		Context("when .properties.enable_log_cache is set to false", func() {
			It("sets the log-cache.disable manifest property to true", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.enable_log_cache": false,
				})
				Expect(err).NotTo(HaveOccurred())

				logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
				Expect(err).NotTo(HaveOccurred())

				disabledProperty, err := logCache.Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabledProperty).To(BeTrue())
			})
		})
	})

	Describe("log cache scheduler", func() {
		var instanceGroup string
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}
		})

		It("has a scheduler", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "log-cache-scheduler")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("syslog forwarding", func() {
		It("includes the vcap rule and does not forward debug logs", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.syslog_host": "example.com",
			})
			Expect(err).NotTo(HaveOccurred())

			syslogForwarder, err := manifest.FindInstanceGroupJob("router", "syslog_forwarder")
			Expect(err).NotTo(HaveOccurred())

			syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
			Expect(err).NotTo(HaveOccurred())
			Expect(syslogConfig).To(ContainSubstring(`if ($programname startswith "vcap.") then stop`))
			Expect(syslogConfig).To(ContainSubstring(`if ($msg contains "DEBUG") then stop`))
		})

		Context("when debug logs are enabled", func() {
			It("does not include the debug stop rule", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.syslog_host":       "example.com",
					".properties.syslog_drop_debug": false,
				})
				Expect(err).NotTo(HaveOccurred())

				syslogForwarder, err := manifest.FindInstanceGroupJob("router", "syslog_forwarder")
				Expect(err).NotTo(HaveOccurred())

				syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
				Expect(err).NotTo(HaveOccurred())
				Expect(syslogConfig).To(ContainSubstring(`if ($programname startswith "vcap.") then stop`))
				Expect(syslogConfig).NotTo(ContainSubstring(`if ($msg contains "DEBUG") then stop`))
			})
		})

		Context("when iptables logs are enabled", func() {
			It("adds a kernel rule", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.syslog_host": "example.com",
					".properties.container_networking_interface_plugin.silk.enable_log_traffic": true,
				})
				Expect(err).NotTo(HaveOccurred())

				syslogForwarder, err := manifest.FindInstanceGroupJob("router", "syslog_forwarder")
				Expect(err).NotTo(HaveOccurred())

				syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
				Expect(err).NotTo(HaveOccurred())
				Expect(syslogConfig).To(ContainSubstring(`if $programname == 'kernel' and ($msg contains "DENY_" or $msg contains "OK_") then -/var/log/kern.log`))
				Expect(syslogConfig).To(ContainSubstring("\n&stop"))
				Expect(syslogConfig).NotTo(ContainSubstring(`"if`)) // previous regression with extra quote
			})
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
					".properties.syslog_host": "example.com",
					".properties.syslog_rule": multilineRule,
				})
				Expect(err).NotTo(HaveOccurred())

				syslogForwarder, err := manifest.FindInstanceGroupJob("router", "syslog_forwarder")
				Expect(err).NotTo(HaveOccurred())

				syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
				Expect(err).NotTo(HaveOccurred())
				Expect(syslogConfig).To(ContainSubstring(`
some
multi
line
rule
`))
			})
		})
	})
})
