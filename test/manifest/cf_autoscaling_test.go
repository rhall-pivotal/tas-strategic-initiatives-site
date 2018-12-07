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

	It("keeps the doc link up-to-date", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
		Expect(err).NotTo(HaveOccurred())

		url, err := job.Property("autoscale/marketplace_documentation_url")
		Expect(err).ToNot(HaveOccurred())
		Expect(url).To(MatchRegexp(`https://docs.pivotal.io/pivotalcf/\d+-\d+/appsman-services/autoscaler/using-autoscaler.html`))
	})

	Context("when the user disables connection pooling", func() {
		It("sets the autoscale api to disable connection pooling", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.autoscale_api_disable_connection_pooling": true,
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
			Expect(err).NotTo(HaveOccurred())

			property, err := job.Property("autoscale/api/disable_connection_pooling")
			Expect(err).ToNot(HaveOccurred())
			Expect(property).To(BeTrue())
		})
	})
})
