package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("windows1803fs", func() {
	It("accepts trusted certs", func() {
		manifest, err := product.RenderManifest(map[string]interface{}{})

		fs, err := manifest.FindInstanceGroupJob("windows_diego_cell", "windows1803fs")
		Expect(err).NotTo(HaveOccurred())

		trustedCerts, err := fs.Property("windows-rootfs/trusted_certs")
		Expect(trustedCerts).NotTo(BeEmpty())
	})
})
