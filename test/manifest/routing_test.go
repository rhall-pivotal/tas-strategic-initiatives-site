package manifest_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Routing", func() {
	Describe("operator defaults", func() {
		It("configures the ha-proxy and router minimum TLS versions", func() {
			manifest, err := product.RenderService.RenderManifest(productConfig)
			Expect(err).NotTo(HaveOccurred())

			haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxy.Property("ha_proxy/disable_tls_10")).To(BeTrue())
			Expect(haproxy.Property("ha_proxy/disable_tls_11")).To(BeTrue())

			router, err := manifest.FindInstanceGroupJob("router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/min_tls_version")).To(Equal("TLSv1.2"))

		})

		XContext("when the operator sets the minimum TLS version to 1.1", func() {
			var (
				manifest planitest.Manifest
				err      error
			)

			BeforeEach(func() {
				if os.Getenv("RENDERER") == "om" {
					err = product.Configure(map[string]interface{}{
						".properties.routing_minimum_tls_version": "tls_v1_1",
					})
					Expect(err).NotTo(HaveOccurred())
				}

				manifest, err = product.RenderService.RenderManifest(productConfig)
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures the ha-proxy and router minimum TLS versions", func() {
				haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy/disable_tls_10")).To(BeTrue())
				Expect(haproxy.Property("ha_proxy/disable_tls_11")).To(BeFalse())

				router, err := manifest.FindInstanceGroupJob("router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/min_tls_version")).To(Equal("TLSv1.1"))
			})
		})

	})
})
