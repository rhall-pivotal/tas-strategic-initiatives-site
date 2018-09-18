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
