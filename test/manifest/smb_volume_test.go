package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SMB volume service", func() {
	Context("when SMB volume services are disabled", func() {
		It("disables the smbdriver job", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.enable_smb_volume_driver": false,
			})
			Expect(err).NotTo(HaveOccurred())

			smbDriver, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "smbdriver")
			Expect(err).NotTo(HaveOccurred())

			smbDriverDisabled, err := smbDriver.Property("disable")
			Expect(err).NotTo(HaveOccurred())

			Expect(smbDriverDisabled).To(BeTrue())
		})
	})

	Context("when SMB volume services are enabled", func() {
		It("enables the smbdriver job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			smbDriver, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "smbdriver")
			Expect(err).NotTo(HaveOccurred())

			smbDriverDisabled, err := smbDriver.Property("disable")
			Expect(err).NotTo(HaveOccurred())

			Expect(smbDriverDisabled).To(BeFalse())
		})
	})
})
