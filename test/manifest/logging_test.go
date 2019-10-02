package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Logging", func() {
	var (
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
		var (
			productTag string
		)

		BeforeEach(func() {
			if productName == "srt" {
				productTag = "Small Footprint Pivotal Application Service"
			} else {
				productTag = "Pivotal Application Service"
			}
		})

		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

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

				grpcPort, err := agent.Property("grpc_port")
				Expect(err).NotTo(HaveOccurred())
				Expect(grpcPort).To(Equal(3459))

				udpDisabled, err := agent.Property("disable_udp")
				Expect(err).NotTo(HaveOccurred())
				Expect(udpDisabled).To(BeTrue())

				By("adding tags to the metrics emitted")
				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(tags).To(HaveKeyWithValue("product", productTag))
				Expect(tags).NotTo(HaveKey("product_version"))
				Expect(tags).To(HaveKeyWithValue("system_domain", "sys.example.com"))
			}
		})
	})

	Describe("system metrics agent", func() {
		It("sets defaults on the system-metrics agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

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

	Describe("system metric scraper", func() {
		var instanceGroup string
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}
		})

		It("configures the system-metric-scraper", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			metricScraper, err := manifest.FindInstanceGroupJob(instanceGroup, "loggr-metric-scraper")
			Expect(err).NotTo(HaveOccurred())

			tlsProps, err := metricScraper.Property("system_metrics/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))

			expectSecureMetrics(metricScraper)
		})

		It("has a leadership-election job collocated", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			le, err := manifest.FindInstanceGroupJob(instanceGroup, "leadership-election")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := le.Property("port")
			Expect(err).ToNot(HaveOccurred())
			Expect(enabled).To(Equal(7100))
		})
	})

	Describe("prom scraper", func() {
		It("configures the prom scraper on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

			for _, ig := range instanceGroups {
				scraper, err := manifest.FindInstanceGroupJob(ig, "prom_scraper")
				Expect(err).NotTo(HaveOccurred())

				expectSecureMetrics(scraper)
			}
		})
	})

	Describe("forwarder agent", func() {
		var (
			productTag string
		)

		BeforeEach(func() {
			if productName == "srt" {
				productTag = "Small Footprint Pivotal Application Service"
			} else {
				productTag = "Pivotal Application Service"
			}
		})

		It("sets defaults on the forwarder agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggr-forwarder-agent")
				Expect(err).NotTo(HaveOccurred())

				port, err := agent.Property("port")
				Expect(err).NotTo(HaveOccurred())
				Expect(port).To(Equal(3458))

				expectSecureMetrics(agent)

				By("adding tags to the metrics emitted")
				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(tags).To(HaveKeyWithValue("product", productTag))
				Expect(tags).NotTo(HaveKey("product_version"))
				Expect(tags).To(HaveKeyWithValue("system_domain", "sys.example.com"))
			}
		})
	})

	Describe("syslog agent", func() {
		It("sets defaults on the syslog agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

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

			instanceGroups := getAllInstanceGroups(manifest)
			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "loggr-syslog-agent")
				Expect(err).NotTo(HaveOccurred())

				_, err = agent.Property("aggregate_drains")
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("syslog binding cache", func() {
		It("sets defaults on the syslog binding cache", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			var instanceGroup string
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			cache, err := manifest.FindInstanceGroupJob(instanceGroup, "loggr-syslog-binding-cache")
			Expect(err).NotTo(HaveOccurred())

			port, err := cache.Property("external_port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(9000))

			expectSecureMetrics(cache)
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

			expectSecureMetrics(logCache)
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

		It("has internal tls certs log-cache-gateway <-> log-cache-cf-auth-proxy", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			gateway, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-gateway")
			Expect(err).NotTo(HaveOccurred())

			gatewayProperties, err := gateway.Path("/properties")
			Expect(err).ToNot(HaveOccurred())

			Expect(gatewayProperties).To(HaveKey("proxy_cert"))
			Expect(gatewayProperties).To(HaveKey("proxy_key"))

			authProxy, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-cf-auth-proxy")
			Expect(err).NotTo(HaveOccurred())

			authProxyProperties, err := authProxy.Path("/properties")
			Expect(err).ToNot(HaveOccurred())

			Expect(authProxyProperties).To(HaveKey("proxy_ca_cert"))
			Expect(authProxyProperties).To(HaveKey("external_cert"))
			Expect(authProxyProperties).To(HaveKey("external_key"))
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

			Expect(routes).To(ContainElement(HaveKeyWithValue("tls_port", 8083)))
			Expect(routes).To(ContainElement(HaveKeyWithValue("server_cert_domain_san", "log-cache.service.cf.internal")))
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

			Expect(proxyPort).To(Equal(8083))

			uaaProperties, err := proxy.Property("uaa")
			Expect(err).ToNot(HaveOccurred())

			Expect(uaaProperties).To(HaveKeyWithValue("client_id", "doppler"))
			Expect(uaaProperties).To(HaveKeyWithValue("internal_addr", "https://uaa.service.cf.internal:8443"))

			Expect(uaaProperties).To(HaveKey("ca_cert"))
			Expect(uaaProperties).To(HaveKey("client_secret"))

		})

		It("has a default max per source", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
			Expect(err).NotTo(HaveOccurred())

			maxPerSource, err := logCache.Property("max_per_source")
			Expect(err).ToNot(HaveOccurred())

			if productName == "srt" {
				Expect(maxPerSource).To(Equal(100000))
			} else {
				Expect(maxPerSource).To(Equal(100000))
			}
		})

		It("has a configurable max per source", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.log_cache_max_per_source": 200000,
			})
			Expect(err).NotTo(HaveOccurred())

			logCache, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache")
			Expect(err).NotTo(HaveOccurred())

			maxPerSource, err := logCache.Property("max_per_source")
			Expect(err).ToNot(HaveOccurred())

			if productName == "srt" {
				Expect(maxPerSource).To(Equal(200000))
			} else {
				Expect(maxPerSource).To(Equal(200000))
			}
		})
	})

	Describe("Traffic Controller", func() {
		var instanceGroup string
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "loggregator_trafficcontroller"
			}
		})

		It("configures TLS for egress", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			trafficcontroller, err := manifest.FindInstanceGroupJob(instanceGroup, "loggregator_trafficcontroller")
			Expect(err).NotTo(HaveOccurred())

			loggrProps, err := trafficcontroller.Property("loggregator")
			Expect(err).ToNot(HaveOccurred())
			Expect(loggrProps).To(HaveKey("outgoing_cert"))
			Expect(loggrProps).To(HaveKey("outgoing_key"))

			routeRegistrar, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
			Expect(err).NotTo(HaveOccurred())

			dopplerRoute, err := routeRegistrar.Property("route_registrar/routes/name=doppler")
			Expect(err).ToNot(HaveOccurred())
			Expect(dopplerRoute).To(HaveKeyWithValue("tls_port", 8081))
			Expect(dopplerRoute).To(HaveKeyWithValue("server_cert_domain_san", "doppler.service.cf.internal"))
			Expect(dopplerRoute).To(HaveKeyWithValue("uris", []interface{}{
				"doppler.sys.example.com",
				"*.doppler.sys.example.com",
			}))
		})

		It("deploys the reverse_log_proxy_gateway", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			gateway, err := manifest.FindInstanceGroupJob(instanceGroup, "reverse_log_proxy_gateway")
			Expect(err).NotTo(HaveOccurred())

			expectSecureMetrics(gateway)

			// test for TLS configuration
			httpConfig, err := gateway.Property("http")
			Expect(err).NotTo(HaveOccurred())
			Expect(httpConfig).To(HaveKeyWithValue("address", "0.0.0.0:8088"))
			Expect(httpConfig).To(HaveKey("cert"))
			Expect(httpConfig).To(HaveKey("key"))

			// test for a subset of properties
			capiAddr, err := gateway.Property("cc/capi_internal_addr")
			Expect(err).NotTo(HaveOccurred())
			Expect(capiAddr).To(Equal("https://cloud-controller-ng.service.cf.internal:9023"))
			uaaAddr, err := gateway.Property("uaa/internal_addr")
			Expect(err).NotTo(HaveOccurred())
			Expect(uaaAddr).To(Equal("https://uaa.service.cf.internal:8443"))

			routeRegistrar, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
			Expect(err).NotTo(HaveOccurred())

			rlpGatewayRoute, err := routeRegistrar.Property("route_registrar/routes/name=rlp-gateway")
			Expect(err).ToNot(HaveOccurred())

			Expect(rlpGatewayRoute).To(HaveKeyWithValue("tls_port", 8088))
			Expect(rlpGatewayRoute).To(HaveKeyWithValue("server_cert_domain_san", "reverse-log-proxy.service.cf.internal"))
			Expect(rlpGatewayRoute).To(HaveKeyWithValue("uris", []interface{}{
				"log-stream.sys.example.com",
				"*.log-stream.sys.example.com",
			}))
		})

		It("is enabled by default", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			trafficController, err := manifest.FindInstanceGroupJob(instanceGroup, "loggregator_trafficcontroller")
			Expect(err).NotTo(HaveOccurred())

			Expect(trafficController.Property("traffic_controller/enabled")).To(BeTrue())
		})

		It("can be disabled", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.enable_v1_firehose": false,
			})
			Expect(err).NotTo(HaveOccurred())

			trafficController, err := manifest.FindInstanceGroupJob(instanceGroup, "loggregator_trafficcontroller")
			Expect(err).NotTo(HaveOccurred())

			Expect(trafficController.Property("traffic_controller/enabled")).To(BeFalse())
		})
	})

	Describe("syslog forwarding", func() {
		It("includes the vcap rule and does not forward debug logs", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.syslog_host": "example.com",
			})
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)
			for _, instanceGroup := range instanceGroups {
				syslogForwarder, err := manifest.FindInstanceGroupJob(instanceGroup, "syslog_forwarder")
				Expect(err).NotTo(HaveOccurred())

				syslogConfig, err := syslogForwarder.Property("syslog/custom_rule")
				Expect(err).NotTo(HaveOccurred())
				Expect(syslogConfig).To(ContainSubstring(`if ($programname startswith "vcap.") then stop`))
				Expect(syslogConfig).To(ContainSubstring(`if ($msg contains "DEBUG") then stop`))
			}
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

func expectSecureMetrics(job planitest.Manifest) {
	metricsProps, err := job.Property("metrics")
	Expect(err).ToNot(HaveOccurred())
	Expect(metricsProps).To(HaveKey("ca_cert"))
	Expect(metricsProps).To(HaveKey("cert"))
	Expect(metricsProps).To(HaveKey("key"))
	Expect(metricsProps).To(HaveKey("server_name"))
}
