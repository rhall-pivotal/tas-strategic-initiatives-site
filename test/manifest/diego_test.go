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

		It("configures the diego bbs job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			bbs, err := manifest.FindInstanceGroupJob(instanceGroup, "bbs")
			Expect(err).NotTo(HaveOccurred())

			By("retrying tasks to be more resilient to temporarily constrained resources")
			maxRetries, err := bbs.Property("tasks/max_retries")
			Expect(err).NotTo(HaveOccurred())
			Expect(maxRetries).To(Equal(3))

			By("configuring TLS to the internal database")
			requireSSL, err := bbs.Property("diego/bbs/sql/require_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(requireSSL).To(BeTrue())

			caCert, err := bbs.Property("diego/bbs/sql/ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).To(Equal("fake-ops-manager-ca-certificate"))
		})
	})

	Context("locket", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "diego_database"
			}
		})

		It("configures the diego locket job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			locket, err := manifest.FindInstanceGroupJob(instanceGroup, "locket")
			Expect(err).NotTo(HaveOccurred())

			requireSSL, err := locket.Property("diego/locket/sql/require_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(requireSSL).To(BeTrue())

			caCert, err := locket.Property("diego/locket/sql/ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).To(Equal("fake-ops-manager-ca-certificate"))
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

		It("co-locates the BPM job with all diego jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, diegoJob := range diegoJobs {
				_, err = manifest.FindInstanceGroupJob(diegoJob.InstanceGroup, "bpm")
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("sets bpm.enabled to true", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, diegoJob := range diegoJobs {
				manifestJob, err := manifest.FindInstanceGroupJob(diegoJob.InstanceGroup, diegoJob.Name)
				Expect(err).NotTo(HaveOccurred())

				bpmEnabled, err := manifestJob.Property("bpm/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(bpmEnabled).To(BeTrue())
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
			manifest, err := product.RenderManifest(nil)
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
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Root file systems", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}
		})

		It("colocates the cflinuxfs2-rootfs-setup job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			setup, err := manifest.FindInstanceGroupJob(instanceGroup, "cflinuxfs2-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err := setup.Property("cflinuxfs2-rootfs/trusted_certs")
			Expect(trustedCerts).NotTo(BeEmpty())
		})

		It("colocates the cflinuxfs3-rootfs-setup job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			setup, err := manifest.FindInstanceGroupJob(instanceGroup, "cflinuxfs3-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err := setup.Property("cflinuxfs3-rootfs/trusted_certs")
			Expect(trustedCerts).NotTo(BeEmpty())
		})

		It("configures the preloaded_rootfses on the rep", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob(instanceGroup, "rep")
			Expect(err).NotTo(HaveOccurred())

			preloadedRootfses, err := rep.Property("diego/rep/preloaded_rootfses")
			Expect(err).NotTo(HaveOccurred())

			Expect(preloadedRootfses).To(ContainElement("cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs.tar"))
			Expect(preloadedRootfses).To(ContainElement("cflinuxfs3:/var/vcap/packages/cflinuxfs3/rootfs.tar"))
		})

		It("ensures the standard root filesystems remain in the layer cache", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			garden, err := manifest.FindInstanceGroupJob(instanceGroup, "garden")
			Expect(err).NotTo(HaveOccurred())

			persistentImageList, err := garden.Property("garden/persistent_image_list")
			Expect(err).NotTo(HaveOccurred())

			Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs2/rootfs.tar"))
			Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs3/rootfs.tar"))
		})
	})

	Context("route integrity", func() {

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}
		})

		It("enables the envoy proxy", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob(instanceGroup, "rep")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := rep.Property("containers/proxy/enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(enabled).To(BeTrue())

			additionalMemory, err := rep.Property("containers/proxy/additional_memory_allocation_mb")
			Expect(err).NotTo(HaveOccurred())
			Expect(additionalMemory).To(Equal(32))
		})

		Context("when strict route integrity is enabled", func() {

			var proxyProperties map[interface{}]interface{}

			BeforeEach(func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.route_integrity": "mutual_tls_verify",
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob(instanceGroup, "rep")
				Expect(err).NotTo(HaveOccurred())

				rawProxyProperties, err := rep.Property("containers/proxy")
				Expect(err).NotTo(HaveOccurred())

				proxyProperties = rawProxyProperties.(map[interface{}]interface{})
			})

			It("enables the proxy", func() {
				Expect(proxyProperties["enabled"]).To(BeTrue())
			})

			It("allocates sufficient RAM for the proxy", func() {
				Expect(proxyProperties["additional_memory_allocation_mb"]).To(Equal(32))
			})

			It("requires and verifies client credentials", func() {
				Expect(proxyProperties["require_and_verify_client_certificates"]).To(BeTrue())
			})

			It("specifies the CA that it trusts", func() {
				Expect(proxyProperties).To(HaveKey("trusted_ca_certificates"))
			})

			It("configures the subject alt name to be verified", func() {
				Expect(proxyProperties["verify_subject_alt_name"]).To(Equal([]interface{}{"gorouter.service.cf.internal"}))
			})

			It("disables direct access to container ports", func() {
				Expect(proxyProperties["enable_unproxied_port_mappings"]).To(BeFalse())
			})
		})
	})

	Context("cflinuxfs2-rootfs and cflinuxfs3-rootfs", func() {

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}
		})

		It("configures the trusted certs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			cflinuxfs2RootfsSetup, err := manifest.FindInstanceGroupJob(instanceGroup, "cflinuxfs2-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err := cflinuxfs2RootfsSetup.Property("cflinuxfs2-rootfs/trusted_certs")
			Expect(err).NotTo(HaveOccurred())
			Expect(trustedCerts).NotTo(BeNil())

			cflinuxfs3RootfsSetup, err := manifest.FindInstanceGroupJob(instanceGroup, "cflinuxfs3-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err = cflinuxfs3RootfsSetup.Property("cflinuxfs3-rootfs/trusted_certs")
			Expect(err).NotTo(HaveOccurred())
			Expect(trustedCerts).NotTo(BeNil())
		})
	})

})
