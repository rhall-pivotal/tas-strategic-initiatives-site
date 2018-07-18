package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego", func() {
	var instanceGroup string

	Context("BBS", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "diego_database"
			}
		})

		It("retries tasks to be more resilient to temporarily constrained resources", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			bbs, err := manifest.FindInstanceGroupJob(instanceGroup, "bbs")
			Expect(err).NotTo(HaveOccurred())

			maxRetries, err := bbs.Property("tasks/max_retries")
			Expect(err).NotTo(HaveOccurred())

			Expect(maxRetries).To(Equal(3))
		})

	})

	Context("SSH Proxy", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "diego_brain"
			}
		})

		It("uses the default UAA URL and port configuration", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			sshProxy, err := manifest.FindInstanceGroupJob(instanceGroup, "ssh_proxy")
			Expect(err).NotTo(HaveOccurred())

			uaaProperties, err := sshProxy.Property("diego/ssh_proxy/uaa")
			Expect(err).NotTo(HaveOccurred())

			Expect(uaaProperties).NotTo(HaveKey("url"))
			Expect(uaaProperties).NotTo(HaveKey("port"))
		})

	})

	Context("Persistence", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}
		})

		It("colocates the nfsv3driver job with the mapfs job from the mapfs-release", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
			Expect(err).NotTo(HaveOccurred())
		})

		Context("Broker errands", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "clock_global"
				}
			})

			It("colocates the broker-reigstrar errand in the appropriate instance group", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "broker-registrar")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
