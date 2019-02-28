package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("loggregator agent", func() {
		It("sets defaults on the loggregator agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "Pivotal Application Service for Windows"))
			Expect(tags).NotTo(HaveKey("product_version"))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		Context("when placement tags are configured by the user", func() {
			It("sets the placement tags on the emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".windows_diego_cell.placement_tags": "tag1,tag2",
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "loggregator_agent_windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})
	})
})
