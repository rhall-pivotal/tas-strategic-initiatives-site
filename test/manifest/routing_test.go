package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego Persistence", func() {
	// TODO: stop skipping once ops-manifest supports testing for credentials
	XDescribe("Gorouter provides client certs in request to Diego cells", func() {
		It("creates a backend cert_chain and private_key", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/backends/cert_chain")).NotTo(BeNil())
			Expect(router.Property("router/backends/private_key")).NotTo(BeNil())
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

	Describe("bpm", func() {
		It("co-locates the BPM job with all routing jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_router", "bpm")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
