package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego Persistence", func() {
	// TODO: stop skipping once ops-manifest supports testing for credentials
	XDescribe("Gorouter provides client certs in request to Diego cells", func() {
		It("creates a backend cert_chain and private_key", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			router, err := manifest.FindInstanceGroupJob("isolated_router", "gorouter")
			Expect(err).NotTo(HaveOccurred())
			Expect(router.Property("router/backends/cert_chain")).NotTo(BeNil())
			Expect(router.Property("router/backends/private_key")).NotTo(BeNil())
		})
	})
})
