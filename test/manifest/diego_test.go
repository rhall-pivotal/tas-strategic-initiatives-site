package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego", func() {

	Describe("Persistence", func() {

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

	Describe("Root file systems", func() {

		It("colocates the cflinuxfs2-rootfs-setup job", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			setup, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "cflinuxfs2-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err := setup.Property("cflinuxfs2-rootfs/trusted_certs")
			Expect(trustedCerts).NotTo(BeEmpty())
		})

		It("configures the preloaded_rootfses on the rep", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())

			preloadedRootfses, err := rep.Property("diego/rep/preloaded_rootfses")
			Expect(err).NotTo(HaveOccurred())

			Expect(preloadedRootfses).To(ContainElement("cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs.tar"))
		})

		Context("when grootfs is disabled", func() {

			It("configures the preloaded_rootfses on the rep", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.enable_grootfs": false,
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
				Expect(err).NotTo(HaveOccurred())

				preloadedRootfses, err := rep.Property("diego/rep/preloaded_rootfses")
				Expect(err).NotTo(HaveOccurred())

				Expect(preloadedRootfses).To(ContainElement("cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs"))
			})

		})

	})

	Describe("Garden", func() {

		It("ensures the standard root filesystem remains in the layer cache", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			garden, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden")
			Expect(err).NotTo(HaveOccurred())

			persistentImageList, err := garden.Property("garden/persistent_image_list")
			Expect(err).NotTo(HaveOccurred())

			Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs2/rootfs.tar"))
		})

		Context("when grootfs is disabled", func() {

			It("ensures the standard root filesystem remains in the layer cache", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.enable_grootfs": false,
				})
				Expect(err).NotTo(HaveOccurred())

				garden, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden")
				Expect(err).NotTo(HaveOccurred())

				persistentImageList, err := garden.Property("garden/persistent_image_list")
				Expect(err).NotTo(HaveOccurred())

				Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs2/rootfs"))
			})

		})

	})

})
