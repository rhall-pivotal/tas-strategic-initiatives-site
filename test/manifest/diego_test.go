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

	Describe("BPM", func() {
		var diegoJobs []Job

		BeforeEach(func() {
			if productName == "srt" {
				diegoJobs = []Job{
					{
						InstanceGroup: "control",
						Name:          "bbs",
					},
					{
						InstanceGroup: "control",
						Name:          "locket",
					},
					{
						InstanceGroup: "control",
						Name:          "auctioneer",
					},
					{
						InstanceGroup: "control",
						Name:          "file_server",
					},
					{
						InstanceGroup: "control",
						Name:          "ssh_proxy",
					},
					{
						InstanceGroup: "compute",
						Name:          "rep",
					},
					{
						InstanceGroup: "compute",
						Name:          "route_emitter",
					},
				}
			} else {
				diegoJobs = []Job{
					{
						InstanceGroup: "diego_database",
						Name:          "bbs",
					},
					{
						InstanceGroup: "diego_database",
						Name:          "locket",
					},
					{
						InstanceGroup: "diego_brain",
						Name:          "auctioneer",
					},
					{
						InstanceGroup: "diego_brain",
						Name:          "file_server",
					},
					{
						InstanceGroup: "diego_brain",
						Name:          "ssh_proxy",
					},
					{
						InstanceGroup: "diego_cell",
						Name:          "rep",
					},
					{
						InstanceGroup: "diego_cell",
						Name:          "route_emitter",
					},
				}
			}
		})

		It("bpm.enabled is set to true", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, diegoJob := range diegoJobs {
				bbs, err := manifest.FindInstanceGroupJob(diegoJob.InstanceGroup, diegoJob.Name)
				Expect(err).NotTo(HaveOccurred())

				bbsBpmEnabled, err := bbs.Property("bpm/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(bbsBpmEnabled).To(BeTrue())
			}
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
	})
})
