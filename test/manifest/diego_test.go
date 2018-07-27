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

	Context("BPM", func() {
		Context("Diego Database jobs", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "diego_database"
				}
			})

			Context("bbs", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					bbs, err := manifest.FindInstanceGroupJob(instanceGroup, "bbs")
					Expect(err).NotTo(HaveOccurred())

					bbsBpmEnabled, err := bbs.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(bbsBpmEnabled).To(BeTrue())
				})
			})

			Context("locket", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					locket, err := manifest.FindInstanceGroupJob(instanceGroup, "locket")
					Expect(err).NotTo(HaveOccurred())

					locketBpmEnabled, err := locket.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(locketBpmEnabled).To(BeTrue())
				})
			})
		})

		Context("Diego Brain jobs", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "diego_brain"
				}
			})

			Context("auctioneer", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					auctioneer, err := manifest.FindInstanceGroupJob(instanceGroup, "auctioneer")
					Expect(err).NotTo(HaveOccurred())

					auctioneerBpmEnabled, err := auctioneer.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(auctioneerBpmEnabled).To(BeTrue())
				})
			})

			Context("file_server", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					fileServer, err := manifest.FindInstanceGroupJob(instanceGroup, "file_server")
					Expect(err).NotTo(HaveOccurred())

					fileServerBpmEnabled, err := fileServer.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(fileServerBpmEnabled).To(BeTrue())
				})
			})

			Context("ssh_proxy", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					sshProxy, err := manifest.FindInstanceGroupJob(instanceGroup, "ssh_proxy")
					Expect(err).NotTo(HaveOccurred())

					sshProxyBpmEnabled, err := sshProxy.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(sshProxyBpmEnabled).To(BeTrue())
				})
			})
		})

		Context("Diego Cell jobs", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "compute"
				} else {
					instanceGroup = "diego_cell"
				}
			})

			Context("rep", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					rep, err := manifest.FindInstanceGroupJob(instanceGroup, "rep")
					Expect(err).NotTo(HaveOccurred())

					repBpmEnabled, err := rep.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(repBpmEnabled).To(BeTrue())
				})
			})

			Context("route_emitter", func() {
				It("bpm.enabled is set to true", func() {
					manifest, err := product.RenderService.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					routeEmitter, err := manifest.FindInstanceGroupJob(instanceGroup, "route_emitter")
					Expect(err).NotTo(HaveOccurred())

					routeEmitterBpmEnabled, err := routeEmitter.Property("bpm/enabled")
					Expect(err).NotTo(HaveOccurred())

					Expect(routeEmitterBpmEnabled).To(BeTrue())
				})
			})
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
