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
				".properties.enable_garden_containerd_mode": false,
			})
			garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

			containerdMode, err := garden.Property("garden/containerd_mode")
			Expect(err).NotTo(HaveOccurred())
			Expect(containerdMode).To(BeFalse())
		})
	})

	Describe("grootfs garbage collection", func() {
		It("sets the reserved disk space", func() {
			manifest := renderProductManifest(product, nil)
			garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

			reservedInMB, err := garden.Property("grootfs/reserved_space_for_other_jobs_in_mb")
			Expect(err).NotTo(HaveOccurred())
			Expect(reservedInMB).To(Equal(15360))
		})

		When("reserved_space_for_other_jobs_in_mb is set", func() {
			It("sets the reserved disk space", func() {
				manifest := renderProductManifest(product, map[string]interface{}{
					".properties.garden_disk_cleanup":                                              "reserved",
					".properties.garden_disk_cleanup.reserved.reserved_space_for_other_jobs_in_mb": 15361,
				})
				garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

				reservedInMB, err := garden.Property("grootfs/reserved_space_for_other_jobs_in_mb")
				Expect(err).NotTo(HaveOccurred())
				Expect(reservedInMB).To(Equal(15361))
			})
		})
	})

	It("ensures the standard root filesystems remain in the layer cache", func() {
		manifest := renderProductManifest(product, nil)
		garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

		persistentImageList, err := garden.Property("garden/persistent_image_list")
		Expect(err).NotTo(HaveOccurred())
		Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs3/rootfs.tar"))
	})

	Describe("log format", func() {
		When("logging_format_timestamp is set to deprecated", func() {
			It("is used in the garden job", func() {
				manifest := renderProductManifest(product, map[string]interface{}{
					".properties.logging_timestamp_format": "deprecated",
				})
				garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

				loggingFormatTimestamp, err := garden.Property("logging/format/timestamp")
				Expect(err).NotTo(HaveOccurred())
				Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))
			})
		})

		When("logging_format_timestamp is set to rfc3339", func() {
			It("the default is used in the garden job", func() {
				manifest := renderProductManifest(product, map[string]interface{}{
					".properties.logging_timestamp_format": "rfc3339",
				})
				garden := findManifestInstanceGroupJob(manifest, instanceGroup, "garden")

				loggingFormatTimestamp, err := garden.Property("logging/format/timestamp")
				Expect(err).NotTo(HaveOccurred())
				Expect(loggingFormatTimestamp).To(Equal("rfc3339"))
			})
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
