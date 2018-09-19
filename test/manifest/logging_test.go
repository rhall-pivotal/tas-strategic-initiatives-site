package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("metron agent", func() {
		It("sets tags on the metron agent", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := []string{
				"isolated_diego_cell",
				"isolated_ha_proxy",
				"isolated_router",
			}

			for _, ig := range instanceGroups {
				agent, err := manifest.FindInstanceGroupJob(ig, "metron_agent")
				Expect(err).NotTo(HaveOccurred())

				tags, err := agent.Property("metron_agent/tags")
				Expect(err).NotTo(HaveOccurred(), "Instance Group: %s", ig)
				Expect(tags).To(HaveKeyWithValue("placement_tag", "isosegtag"))
				Expect(tags).To(HaveKeyWithValue("product", "PCF Isolation Segment"))
				Expect(tags).To(HaveKeyWithValue("product_version", MatchRegexp(`^\d+\.\d+\.\d+.*`)))
				Expect(tags).To(HaveKeyWithValue("system_domain", Not(BeEmpty())))
			}
		})
	})
})
