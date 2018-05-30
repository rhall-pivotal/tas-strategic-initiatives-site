package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego Persistence", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "compute"
		} else {
			instanceGroup = "diego_cell"
		}
	})

	It("colocates the nfsv3driver job with the mapfs job from the mapfs-release", func() {
		manifest, err := product.RenderService.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		_, err = manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
		Expect(err).NotTo(HaveOccurred())

		_, err = manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
		Expect(err).NotTo(HaveOccurred())
	})
})
