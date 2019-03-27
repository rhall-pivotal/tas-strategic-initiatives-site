package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Garden", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "compute"
		} else {
			instanceGroup = "diego_cell"
		}
	})

	It("enables containerd_mode by default", func() {
		manifest := renderProductManifest(product, nil)
		garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

		containerdMode, err := garden.Property("garden/containerd_mode")
		Expect(err).NotTo(HaveOccurred())
		Expect(containerdMode).To(BeTrue())
	})

	When("opted out of containerd mode", func() {
		It("disables containerd_mode", func() {
			manifest := renderProductManifest(product, map[string]interface{}{
				".properties.garden_containerd_opt_out": true,
			})
			garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

			containerdMode, err := garden.Property("garden/containerd_mode")
			Expect(err).NotTo(HaveOccurred())
			Expect(containerdMode).To(BeFalse())
		})
	})
})

func renderProductManifest(p *planitest.ProductService, c map[string]interface{}) planitest.Manifest {
	manifest, err := p.RenderManifest(c)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return manifest
}

func findManifestInstanceGroupJob(m planitest.Manifest, group, job string) planitest.Manifest {
	manifest, err := m.FindInstanceGroupJob(group, job)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return manifest
}
