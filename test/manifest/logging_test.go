package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	var instanceGroup string

	Describe("traffic controller", func() {

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "loggregator_trafficcontroller"
			}
		})

		It("disables support for forwarding syslog to metron", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob(instanceGroup, "metron_agent")
			Expect(err).NotTo(HaveOccurred())

			syslogForwardingEnabled, err := agent.Property("syslog_daemon_config/enable")
			Expect(syslogForwardingEnabled).To(BeFalse())
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

		It("has tls server certs", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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

		It("has a log-cache-group-reader with a port", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			groupReader, err := manifest.FindInstanceGroupJob(instanceGroup, "log-cache-group-reader")
			Expect(err).NotTo(HaveOccurred())

			port, err := groupReader.Property("port")
			Expect(err).ToNot(HaveOccurred())
			if productName == "srt" {
				Expect(port).To(Equal(8088))
			} else {
				Expect(port).To(Equal(8084))
			}
		})

		It("has a log-cache-nozzle with tls certs", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
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
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "log-cache-scheduler")
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
