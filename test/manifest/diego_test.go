package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rep Windows", func() {
	It("sets memory and disk capacity", func() {
		manifest, err := product.RenderManifest(map[string]interface{}{
			".windows_diego_cell.executor_memory_capacity": "500",
			".windows_diego_cell.executor_disk_capacity":   "500",
		})
		Expect(err).NotTo(HaveOccurred())

		repWindows, err := manifest.FindInstanceGroupJob("windows_diego_cell", "rep_windows")
		Expect(err).NotTo(HaveOccurred())

		azureFaultDomains, err := repWindows.Property("diego/rep/use_azure_fault_domains")
		Expect(err).NotTo(HaveOccurred())
		Expect(azureFaultDomains).To(BeTrue())

		memoryCapacity, err := repWindows.Property("diego/executor/memory_capacity_mb")
		Expect(err).NotTo(HaveOccurred())
		Expect(memoryCapacity).To(Equal(500))

		diskCapacity, err := repWindows.Property("diego/executor/disk_capacity_mb")
		Expect(err).NotTo(HaveOccurred())
		Expect(diskCapacity).To(Equal(500))
	})

	Context("instance identity", func() {
		It("uses an intermediate CA cert from Credhub", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("windows_diego_cell", "rep_windows")
			Expect(err).NotTo(HaveOccurred())

			caCert, err := rep.Property("diego/executor/instance_identity_ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).To(Equal("((diego-instance-identity-intermediate-ca-2018.certificate))"))

			caKey, err := rep.Property("diego/executor/instance_identity_key")
			Expect(err).NotTo(HaveOccurred())
			Expect(caKey).To(Equal("((diego-instance-identity-intermediate-ca-2018.private_key))"))
		})
	})
})
