package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Routing", func() {
	Describe("router_headers_remove_if_specified", func() {
		var (
			inputProperties     map[string]interface{}
			routerInstanceGroup string
		)

		BeforeEach(func() {
			routerInstanceGroup = "isolated_router"
			inputProperties = map[string]interface{}{
				".properties.router_headers_remove_if_specified": []map[string]interface{}{
					{
						"name": "header1",
					},
					{
						"name": "header2",
					},
				}}
		})

		It("sets the headers to be removed for http responses", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
			Expect(err).NotTo(HaveOccurred())

			removeHeaders, err := job.Property("router/http_rewrite/responses/remove_headers")
			Expect(err).NotTo(HaveOccurred())
			Expect(removeHeaders.([]interface{})[0].(map[interface{}]interface{})["name"]).To(Equal("header1"))
			Expect(removeHeaders.([]interface{})[1].(map[interface{}]interface{})["name"]).To(Equal("header2"))
		})
	})

	Describe("Gorouter provides client certs in request to Diego cells", func() {

		It("creates a backend cert_chain and private_key", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())

			_, err = router.Property("router/backends/cert_chain")
			Expect(err).NotTo(HaveOccurred())

			_, err = router.Property("router/backends/private_key")
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("idle timeouts", func() {

		It("inherits the PAS frontend idle timeout", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			haproxy, err := manifest.FindInstanceGroupJob("isolated_ha_proxy", "haproxy")
			Expect(err).NotTo(HaveOccurred())
			haproxyTimeout, err := haproxy.Property("ha_proxy/keepalive_timeout")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxyTimeout).To(Equal(900))

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			routerTimeout, err := router.Property("router/frontend_idle_timeout")
			Expect(err).NotTo(HaveOccurred())
			Expect(routerTimeout).To(Equal(900))
		})

	})

	Describe("logging", func() {
		It("sets defaults on the udp forwarder for the router", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_router", "loggr-udp-forwarder")
			Expect(err).NotTo(HaveOccurred())

			udpForwarder, err := manifest.FindInstanceGroupJob("isolated_router", "loggr-udp-forwarder")
			Expect(err).NotTo(HaveOccurred())

			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("ca"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("cert"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("key"))
		})
	})

	Describe("bpm", func() {
		It("co-locates the BPM job with all routing jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_router", "bpm")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Route Services", func() {
		It("disables route services internal lookup when internal_lookup is false", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.route_services_internal_lookup": false,
			})
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/route_services_internal_lookup")).To(Equal(false))
		})

		It("enables route services internal lookup when internal_lookup is true", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.route_services_internal_lookup": true,
			})
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/route_services_internal_lookup")).To(Equal(true))
		})
	})

	Describe("Route Balancer", func() {
		It("set balancing_algorithm to the value of router_balancing_algorithm property", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.router_balancing_algorithm": "least-connection",
			})
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/balancing_algorithm")).To(Equal("least-connection"))
		})
	})

	Describe("isolation_segments", func() {
		Context("when compute isolation is enabled", func() {
			It("adds the appropriate placement_tag", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
				Expect(err).NotTo(HaveOccurred())

				placementTag, err := router.Property("router/isolation_segments")
				Expect(err).NotTo(HaveOccurred())
				Expect(placementTag).To(ContainElement("isosegtag"))
			})
		})

		Context("when compute isolation is disabled", func() {
			It("does not have a placement", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.compute_isolation":                                "disabled",
					".properties.compute_isolation.enabled.isolation_segment_name": "",
				})
				Expect(err).NotTo(HaveOccurred())

				router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
				Expect(err).NotTo(HaveOccurred())

				placementTag, err := router.Property("router/isolation_segments")
				Expect(err).NotTo(HaveOccurred())
				Expect(placementTag).To(BeEmpty())
			})
		})
	})

	Describe("services ca", func() {
		It("adds the /services/intermediate_tls_ca to the router ca_certs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())

			routerCACerts, err := router.Property("router/ca_certs")
			Expect(err).NotTo(HaveOccurred())
			Expect(routerCACerts).NotTo(BeEmpty())
			Expect(routerCACerts).To(ContainSubstring("((/services/intermediate_tls_ca.ca))"))
		})
	})
})
