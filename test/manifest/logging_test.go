package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	var instanceGroup string

	FDescribe("traffic controller", func() {

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "loggregator_trafficcontroller"
			}
		})

		It("sets defaults on the metron agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob(instanceGroup, "metron_agent")
			Expect(err).NotTo(HaveOccurred())

			By("disabling support for forwarding syslog to metron")
			syslogForwardingEnabled, err := agent.Property("syslog_daemon_config/enable")
			Expect(err).NotTo(HaveOccurred())
			Expect(syslogForwardingEnabled).To(BeFalse())

			By("disabling the cf deployment name in emitted metrics")
			deploymentName, err := agent.Property("metron_agent/deployment")
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentName).To(Equal(""))
		})

		Context("when upgrading", func() {
			It("enables the cf deployment name in emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.enable_cf_metric_name": true,
				})

				agent, err := manifest.FindInstanceGroupJob(instanceGroup, "metron_agent")
				Expect(err).NotTo(HaveOccurred())

				deploymentName, err := agent.Property("metron_agent/deployment")
				Expect(err).NotTo(HaveOccurred())
				Expect(deploymentName).To(Equal("cf"))
			})
		})
	})

	Describe("log cache", func() {
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
	})
})
