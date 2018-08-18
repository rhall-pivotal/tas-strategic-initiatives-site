package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("CF CLI", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
	})

	It("colocates the cf-cli-6-linux job on the instance group used to run errands", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		_, err = manifest.FindInstanceGroupJob(instanceGroup, "cf-cli-6-linux")
		Expect(err).NotTo(HaveOccurred())
	})
})
