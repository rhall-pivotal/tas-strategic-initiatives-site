package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("DNS search domain", func() {
		It("configures search_domains on the garden-cni job", func() {
			inputProperties := map[string]interface{}{
				".properties.cf_networking_search_domains": "some-search-domain,another-search-domain",
			}

			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(Equal([]interface{}{
				"some-search-domain",
				"another-search-domain",
			}))
		})

		It("defaults search_domains to empty", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(HaveLen(0))
		})
	})

	Describe("Network Policy Enforcement", func() {
		It("defaults the rendered manifest job property to false", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())
			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "vxlan-policy-agent")
			Expect(err).NotTo(HaveOccurred())
			disabled, err := job.Property("disable_container_network_policy")
			Expect(err).NotTo(HaveOccurred())
			Expect(disabled).To(BeFalse())
		})

		It("configures the vxlan-policy-agent job when setting it disabled in the UI", func() {
			inputProperties := map[string]interface{}{
				".properties.enable_silk_policy_enforcement": false,
			}
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())
			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "vxlan-policy-agent")
			Expect(err).NotTo(HaveOccurred())
			disabled, err := job.Property("disable_container_network_policy")
			Expect(err).NotTo(HaveOccurred())
			Expect(disabled).To(BeTrue())
		})

		It("configures the vxlan-policy-agent job when setting it disabled in the UI", func() {
			inputProperties := map[string]interface{}{
				".properties.enable_silk_policy_enforcement": false,
			}
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())
			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "vxlan-policy-agent")
			Expect(err).NotTo(HaveOccurred())
			disabled, err := job.Property("disable_container_network_policy")
			Expect(err).NotTo(HaveOccurred())
			Expect(disabled).To(BeTrue())
		})
	})

	Describe("BOSH DNS Adapter for App Service Discovery", func() {
		It("colocates the dns-adapter and route emitter on the isolated_diego_cell", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_diego_cell", "bosh-dns-adapter")
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "route_emitter")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := job.Property("internal_routes/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(enabled).To(BeTrue())
		})

		Context("when PAS internal domain is empty", func() {
			It("defaults internal domain to apps.internal", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "bosh-dns-adapter")
				Expect(err).NotTo(HaveOccurred())

				internalDomains, err := job.Property("internal_domains")
				Expect(err).NotTo(HaveOccurred())

				Expect(internalDomains).To(ConsistOf("apps.internal"))
			})
		})
	})

	Describe("TLS termination", func() {
		Context("when TLS is terminated for the first time at infrastructure load balancer", func() {
			It("configures the router and proxy", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("isolated_ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy")).ShouldNot(HaveKey("client_ca_file"))
				Expect(haproxy.Property("ha_proxy/client_cert")).To(BeFalse())

				router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/forwarded_client_cert")).To(ContainSubstring("always_forward"))
			})
		})

		Context("when TLS is terminated for the first time at ha proxy", func() {
			Context("when ha proxy client cert validation is set to none", func() {
				It("configures ha proxy and router", func() {
					manifest, err := product.RenderManifest(map[string]interface{}{
						".properties.routing_tls_termination": "ha_proxy",
					})
					Expect(err).NotTo(HaveOccurred())

					haproxy, err := manifest.FindInstanceGroupJob("isolated_ha_proxy", "haproxy")
					Expect(err).NotTo(HaveOccurred())
					Expect(haproxy.Property("ha_proxy")).ShouldNot(HaveKey("client_ca_file"))
					Expect(haproxy.Property("ha_proxy/client_cert")).To(BeFalse())

					router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
					Expect(err).NotTo(HaveOccurred())
					Expect(router.Property("router/forwarded_client_cert")).To(ContainSubstring("forward"))
				})
			})

			Context("when ha proxy client cert validation is set to request ", func() {
				It("configures ha proxy and router", func() {
					manifest, err := product.RenderManifest(map[string]interface{}{
						".properties.routing_tls_termination":        "ha_proxy",
						".properties.haproxy_client_cert_validation": "request",
					})
					Expect(err).NotTo(HaveOccurred())

					haproxy, err := manifest.FindInstanceGroupJob("isolated_ha_proxy", "haproxy")
					Expect(err).NotTo(HaveOccurred())
					Expect(haproxy.Property("ha_proxy")).ShouldNot(HaveKey("client_ca_file"))
					Expect(haproxy.Property("ha_proxy/client_cert")).To(BeTrue())

					router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
					Expect(err).NotTo(HaveOccurred())
					Expect(router.Property("router/forwarded_client_cert")).To(ContainSubstring("forward"))
				})
			})
		})

		Context("when TLS is terminated for the first time at the router", func() {
			It("configures the router and proxy", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.routing_tls_termination": "router",
				})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("isolated_ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy")).ShouldNot(HaveKey("client_ca_file"))
				Expect(haproxy.Property("ha_proxy/client_cert")).To(BeFalse())

				router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/forwarded_client_cert")).To(ContainSubstring("sanitize_set"))
			})
		})
	})

	Describe("Routing", func() {
		Describe("drain_timeout", func() {
			var (
				inputProperties     map[string]interface{}
				routerInstanceGroup string
			)

			BeforeEach(func() {
				routerInstanceGroup = "isolated_router"
			})

			Describe("when the property is set", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".isolated_router.drain_timeout": 999,
					}
				})

				It("sets the drain_timeout", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
					Expect(err).NotTo(HaveOccurred())

					drainTimeout, err := job.Property("router/drain_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(drainTimeout).To(Equal(999))
				})
			})

			Describe("when the property is not set", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{}
				})

				It("defaults to false", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
					Expect(err).NotTo(HaveOccurred())

					drainTimeout, err := job.Property("router/drain_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(drainTimeout).To(Equal(900))
				})
			})
		})
	})
})
