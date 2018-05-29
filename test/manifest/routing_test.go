package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Routing", func() {
	Describe("operator defaults", func() {
		It("configures the ha-proxy and router minimum TLS versions", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxy.Property("ha_proxy/disable_tls_10")).To(BeTrue())
			Expect(haproxy.Property("ha_proxy/disable_tls_11")).To(BeTrue())

			router, err := manifest.FindInstanceGroupJob("router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/min_tls_version")).To(Equal("TLSv1.2"))

		})

		Context("when the operator sets the minimum TLS version to 1.1", func() {
			var (
				manifest planitest.Manifest
				err      error
			)

			BeforeEach(func() {
				manifest, err = product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_minimum_tls_version": "tls_v1_1",
				})
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

	Describe("IP Logging", func() {
		Context("when the operator chooses to log client Ips", func() {
			It("does not disable ip logging or x-forwarded-for logging", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_log_client_ips": "log_client_ips",
				})
				Expect(err).NotTo(HaveOccurred())

				router, err := manifest.FindInstanceGroupJob("router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/disable_log_forwarded_for")).To(BeFalse())
				Expect(router.Property("router/disable_log_source_ips")).To(BeFalse())
			})
		})
		Context("when the operator chooses `Disable logging of X-Forwarded-For header only`", func() {
			It("only disables x-forwarded-for logging but not source ip logging", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_log_client_ips": "disable_x_forwarded_for",
				})
				Expect(err).NotTo(HaveOccurred())

				router, err := manifest.FindInstanceGroupJob("router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/disable_log_forwarded_for")).To(BeTrue())
				Expect(router.Property("router/disable_log_source_ips")).To(BeFalse())
			})
		})
		Context("when the operator chooses `Disable logging of both source IP and X-Forwarded-For header`", func() {
			It("disbales both source ip logging and x-forwarded-for logging", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_log_client_ips": "disable_all_log_client_ips",
				})
				Expect(err).NotTo(HaveOccurred())

				router, err := manifest.FindInstanceGroupJob("router", "gorouter")
				Expect(err).NotTo(HaveOccurred())
				Expect(router.Property("router/disable_log_forwarded_for")).To(BeTrue())
				Expect(router.Property("router/disable_log_source_ips")).To(BeTrue())
			})
		})
	})
})
