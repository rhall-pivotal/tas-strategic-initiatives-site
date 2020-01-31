package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("windows2019fs", func() {
	It("accepts trusted certs", func() {
		manifest, err := product.RenderManifest(map[string]interface{}{})

		fs, err := manifest.FindInstanceGroupJob("windows_diego_cell", "windows2019fs")
		Expect(err).NotTo(HaveOccurred())

		trustedCerts, err := fs.Property("windows-rootfs/trusted_certs")
		Expect(trustedCerts).NotTo(BeEmpty())
		Expect(trustedCerts).To(ContainSubstring("((/services/intermediate_tls_ca.ca))"))
	})
})
