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

			haproxyDisableTLS10, err := haproxy.Property("ha_proxy/disable_tls_10")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxyDisableTLS10).To(BeTrue())

			haproxyDisableTLS11, err := haproxy.Property("ha_proxy/disable_tls_11")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxyDisableTLS11).To(BeTrue())

			router, err := manifest.FindInstanceGroupJob("router", "gorouter")
			Expect(err).NotTo(HaveOccurred())

			routerMinTLSVersion, err := router.Property("router/min_tls_version")
			Expect(err).NotTo(HaveOccurred())
			Expect(routerMinTLSVersion).To(Equal("TLSv1.2"))
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

				haproxyDisableTLS10, err := haproxy.Property("ha_proxy/disable_tls_10")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxyDisableTLS10).To(BeTrue())

				haproxyDisableTLS11, err := haproxy.Property("ha_proxy/disable_tls_11")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxyDisableTLS11).To(BeFalse())

				router, err := manifest.FindInstanceGroupJob("router", "gorouter")
				Expect(err).NotTo(HaveOccurred())

				routerMinTLSVersion, err := router.Property("router/min_tls_version")
				Expect(err).NotTo(HaveOccurred())
				Expect(routerMinTLSVersion).To(Equal("TLSv1.1"))
			})

		})
	})

	Describe("TLS termination", func() {

		It("secures traffic between the infrastructure load balancer and HAProxy / Gorouter", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
			Expect(err).NotTo(HaveOccurred())

			haproxySSLPEM, err := haproxy.Property("ha_proxy/ssl_pem")
			Expect(err).NotTo(HaveOccurred())
			Expect(haproxySSLPEM).NotTo(BeEmpty())

			router, err := manifest.FindInstanceGroupJob("router", "gorouter")
			Expect(err).NotTo(HaveOccurred())

			routerEnableSSL, err := router.Property("router/enable_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(routerEnableSSL).To(BeTrue())

			routerTLSPEM, err := router.Property("router/tls_pem")
			Expect(err).NotTo(HaveOccurred())
			Expect(routerTLSPEM).NotTo(BeEmpty())
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

				disableLogForwardedFor, err := router.Property("router/disable_log_forwarded_for")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogForwardedFor).To(BeFalse())

				disableLogSourceIPs, err := router.Property("router/disable_log_source_ips")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogSourceIPs).To(BeFalse())
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

				disableLogForwardedFor, err := router.Property("router/disable_log_forwarded_for")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogForwardedFor).To(BeTrue())

				disableLogSourceIPs, err := router.Property("router/disable_log_source_ips")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogSourceIPs).To(BeFalse())
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

				disableLogForwardedFor, err := router.Property("router/disable_log_forwarded_for")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogForwardedFor).To(BeTrue())

				disableLogSourceIPs, err := router.Property("router/disable_log_source_ips")
				Expect(err).NotTo(HaveOccurred())
				Expect(disableLogSourceIPs).To(BeTrue())
			})
		})
	})

	// TODO: stop skipping once ops-manifest supports testing for credentials
	XDescribe("Gorouter provides client certs in request to Diego cells", func() {
		It("creates a backend cert_chain and private_key", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("router", "gorouter")
			Expect(err).NotTo(HaveOccurred())

			certChain, err := router.Property("router/backends/cert_chain")
			Expect(err).NotTo(HaveOccurred())
			Expect(certChain).NotTo(BeNil())

			privateKey, err := router.Property("router/backends/private_key")
			Expect(err).NotTo(HaveOccurred())
			Expect(privateKey).NotTo(BeNil())
		})
	})

	Describe("TLS termination", func() {
		Context("when TLS is terminated for the first time at infrastructure load balancer", func() {
			It("sets ha_proxy.client_ca_file", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy/client_ca_file")).NotTo(BeNil())
			})

			It("sets ha_proxy.client_cert to false", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy/client_cert")).To(BeFalse())
			})
		})

		Context("when TLS is terminated for the first time at ha proxy", func() {
			It("sets ha_proxy.client_ca_file", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_tls_termination": "ha_proxy",
				})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy/client_ca_file")).NotTo(BeNil())
			})

			It("sets ha_proxy.client_cert to true", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.routing_tls_termination": "ha_proxy",
				})
				Expect(err).NotTo(HaveOccurred())

				haproxy, err := manifest.FindInstanceGroupJob("ha_proxy", "haproxy")
				Expect(err).NotTo(HaveOccurred())
				Expect(haproxy.Property("ha_proxy/client_cert")).To(BeTrue())
			})
		})
	})
})
