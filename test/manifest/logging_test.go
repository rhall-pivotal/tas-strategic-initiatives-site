package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("metron agent", func() {
		It("sets defaults on the metron agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metron_agent_windows")
			Expect(err).NotTo(HaveOccurred())

			By("disabling the cf deployment name in emitted metrics")
			deploymentName, err := agent.Property("metron_agent/deployment")
			Expect(err).NotTo(HaveOccurred())
			Expect(deploymentName).To(Equal(""))

			By("setting tags on the emitted metrics")
			tags, err := agent.Property("metron_agent/tags")
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveKeyWithValue("product", "Pivotal Application Service for Windows"))
			Expect(tags).To(HaveKeyWithValue("product_version", MatchRegexp(`^\d+\.\d+\.\d+.*`)))
			Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
		})

		Context("when placement tags are configured by the user", func() {
			It("sets the placement tags on the emitted metrics", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".windows_diego_cell.placement_tags": "tag1,tag2",
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metron_agent_windows")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("metron_agent/tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("placement_tag", "tag1,tag2"))
			})
		})

		Context("when the enable cf metric name is set to true (migration during upgrades)", func() {
			It("sets the metric deployment name to cf", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.enable_cf_metric_name": true,
				})
				Expect(err).NotTo(HaveOccurred())

				agent, err := manifest.FindInstanceGroupJob("windows_diego_cell", "metron_agent_windows")
				Expect(err).NotTo(HaveOccurred())

				deploymentName, err := agent.Property("metron_agent/deployment")
				Expect(err).NotTo(HaveOccurred())
				Expect(deploymentName).To(Equal("cf"))
			})
		})
	})
})
