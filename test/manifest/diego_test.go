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
			Expect(caCert).To(Equal("((diego-instance-identity-intermediate-ca-2-7.certificate))"))

			caKey, err := rep.Property("diego/executor/instance_identity_key")
			Expect(err).NotTo(HaveOccurred())
			Expect(caKey).To(Equal("((diego-instance-identity-intermediate-ca-2-7.private_key))"))
		})
	})

	Describe("logging", func() {
		It("sets defaults on the udp forwarder for the router", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			udpForwarder, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "loggr-udp-forwarder")
			Expect(err).NotTo(HaveOccurred())

			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("ca"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("cert"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("key"))
		})
	})

	//TODO: Testing inheritance from PAS requires manual additions to ops-manifest fixture.
	// Unpend this test when we can render the manifest with inheritance properties like
	// `..cf.properties.cf_networking_interface_plugin`.
	PDescribe("connecting to instance address", func() {
		Context("when container networking plugin is external", func() {
			It("sets advertise_preference_for_instance_address to true", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					"..cf.properties.container_networking_interface_plugin": "external",
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
				Expect(err).NotTo(HaveOccurred())

				property, err := rep.Property("diego/rep/advertise_preference_for_instance_address")
				Expect(err).NotTo(HaveOccurred())

				Expect(property).To(BeTrue())
			})
		})

		Context("when container networking plugin is silk", func() {
			It("sets advertise_preference_for_instance_address to false", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					"..cf.properties.container_networking_interface_plugin": "silk",
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
				Expect(err).NotTo(HaveOccurred())

				property, err := rep.Property("diego/rep/advertise_preference_for_instance_address")
				Expect(err).NotTo(HaveOccurred())

				Expect(property).To(BeFalse())
			})
		})
	})

	Context("Metrics", func() {
		It("sets cpu weight on", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())
			rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
			Expect(err).NotTo(HaveOccurred())
			setCPUWeight, err := rep.Property("containers/set_cpu_weight")
			Expect(err).NotTo(HaveOccurred())
			Expect(setCPUWeight).To(BeTrue())
		})
	})

	Describe("placement_tag", func() {
		Context("when compute isolation is enabled", func() {
			It("adds the appropriate placement_tag", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
				Expect(err).NotTo(HaveOccurred())

				placementTag, err := rep.Property("diego/rep/placement_tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(placementTag).To(ContainElement("isosegtag"))
			})
		})

		Context("when compute isolation is disabled", func() {
			It("does not have a placement", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.compute_isolation":                                "disabled",
					".properties.compute_isolation.enabled.isolation_segment_name": "",
				})
				Expect(err).NotTo(HaveOccurred())

				rep, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "rep")
				Expect(err).NotTo(HaveOccurred())

				properties, err := rep.Property("diego/rep")
				Expect(err).ToNot(HaveOccurred())
				Expect(properties).ToNot(HaveKey("placement_tags"))
			})
		})
	})
})
