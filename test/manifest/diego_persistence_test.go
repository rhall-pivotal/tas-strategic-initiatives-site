package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego Persistence", func() {
	It("colocates the nfsv3driver job with the mapfs job from the mapfs-release", func() {
		instanceGroup := "isolated_diego_cell"

		manifest, err := product.RenderService.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		_, err = manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
		Expect(err).NotTo(HaveOccurred())

		_, err = manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
		Expect(err).NotTo(HaveOccurred())
	})
})
