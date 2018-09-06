package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CF Autoscaling", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
	})

	It("sets the organization and space for the test-autoscale errand", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "test-autoscaling")
		Expect(err).NotTo(HaveOccurred())

		space, err := job.Property("autoscale/space")
		Expect(err).ToNot(HaveOccurred())
		Expect(space).To(Equal("autoscaling"))

		org, err := job.Property("autoscale/organization")
		Expect(err).ToNot(HaveOccurred())
		Expect(org).To(Equal("system"))
	})
})
