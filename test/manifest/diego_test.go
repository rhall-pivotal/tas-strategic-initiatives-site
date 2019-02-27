package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Diego", func() {

	Describe("has BPM enabled", func() {
		It("co-locates bpm job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_diego_cell", "bpm")
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets bpm.enabled to true for rep and route_emitter", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())

			repBpmEnabled, err := rep.Property("bpm/enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(repBpmEnabled).To(BeTrue())

			routeEmitter, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "route_emitter")
			Expect(err).NotTo(HaveOccurred())

			routeEmitterBpmEnabled, err := routeEmitter.Property("bpm/enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(routeEmitterBpmEnabled).To(BeTrue())
		})
	})

	Describe("Persistence", func() {
		It("colocates the nfsv3driver job with the mapfs job from the mapfs-release", func() {
			instanceGroup := "isolated_diego_cell"

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Describe("Root file systems", func() {

		It("colocates the cflinuxfs3-rootfs-setup job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			setup, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "cflinuxfs3-rootfs-setup")
			Expect(err).NotTo(HaveOccurred())

			trustedCerts, err := setup.Property("cflinuxfs3-rootfs/trusted_certs")
			Expect(trustedCerts).NotTo(BeEmpty())
		})

		It("configures the preloaded_rootfses on the rep", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())

			preloadedRootfses, err := rep.Property("diego/rep/preloaded_rootfses")
			Expect(err).NotTo(HaveOccurred())

			Expect(preloadedRootfses).To(ContainElement("cflinuxfs3:/var/vcap/packages/cflinuxfs3/rootfs.tar"))
		})
	})

	Describe("Garden", func() {

		It("ensures the standard root filesystem remains in the layer cache", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			garden, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden")
			Expect(err).NotTo(HaveOccurred())

			persistentImageList, err := garden.Property("garden/persistent_image_list")
			Expect(err).NotTo(HaveOccurred())

			Expect(persistentImageList).To(ContainElement("/var/vcap/packages/cflinuxfs3/rootfs.tar"))
		})
	})

	Context("route integrity", func() {

		var proxyProperties map[interface{}]interface{}

		It("enables the envoy proxy", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())

			rawProxyProperties, err := rep.Property("containers/proxy")
			Expect(err).NotTo(HaveOccurred())

			proxyProperties = rawProxyProperties.(map[interface{}]interface{})

			Expect(proxyProperties["enabled"]).To(BeTrue())
			Expect(proxyProperties["additional_memory_allocation_mb"]).To(Equal(32))
			Expect(proxyProperties).NotTo(HaveKey("enable_unproxied_port_mappings"))
			Expect(proxyProperties).NotTo(HaveKey("require_and_verify_client_certificates"))
			Expect(proxyProperties).NotTo(HaveKey("trusted_ca_certificates"))
			Expect(proxyProperties).NotTo(HaveKey("verify_subject_alt_name"))
		})

		Context("when strict route integrity is enabled", func() {

			BeforeEach(func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.route_integrity": "mutual_tls_verify",
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
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
				Expect(proxyProperties["verify_subject_alt_name"]).To(Equal([]interface{}{
					"gorouter.service.cf.internal",
					"ssh-proxy.service.cf.internal",
				}))
			})

			It("disables direct access to container ports", func() {
				Expect(proxyProperties["enable_unproxied_port_mappings"]).To(BeFalse())
			})
		})
	})

	Context("instance identity", func() {
		It("uses an intermediate CA cert from Credhub", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())

			caCert, err := rep.Property("diego/executor/instance_identity_ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).To(Equal("((diego-instance-identity-intermediate-ca-2018.certificate))"))

			caKey, err := rep.Property("diego/executor/instance_identity_key")
			Expect(err).NotTo(HaveOccurred())
			Expect(caKey).To(Equal("((diego-instance-identity-intermediate-ca-2018.private_key))"))
		})
	})

	Context("garden grootfs garbage collection", func() {
		It("sets the reserved disk space", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			garden, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden")
			Expect(err).NotTo(HaveOccurred())

			reservedInMB, err := garden.Property("grootfs/reserved_space_for_other_jobs_in_mb")
			Expect(err).NotTo(HaveOccurred())
			Expect(reservedInMB).To(Equal(15360))
		})
	})
})
