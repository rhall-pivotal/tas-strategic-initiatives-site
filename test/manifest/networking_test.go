package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("Container networking", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "diego_cell"
			} else {
				instanceGroup = "compute"
			}
		})

		Context("when the operator configures database connection timeout for CNI plugin", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.cf_networking_database_connection_timeout": 250,
				}

				if productName == "ert" {
					instanceGroup = "diego_database"
				} else {
					instanceGroup = "control"
				}
			})

			It("sets the manifest database connection timeout properties for the cf networking jobs to be 250", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				policyServerJob, err := manifest.FindInstanceGroupJob(instanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				policyServerConnectTimeoutSeconds, err := policyServerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerConnectTimeoutSeconds).To(Equal(250))

				policyServerInternalJob, err := manifest.FindInstanceGroupJob(instanceGroup, "policy-server-internal")
				Expect(err).NotTo(HaveOccurred())

				policyServerInternalConnectTimeoutSeconds, err := policyServerInternalJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerInternalConnectTimeoutSeconds).To(Equal(250))

				silkControllerJob, err := manifest.FindInstanceGroupJob(instanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				silkControllerConnectTimeoutSeconds, err := silkControllerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(silkControllerConnectTimeoutSeconds).To(Equal(250))
			})
		})

		Context("when Silk is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "silk",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/silk-cni/config/cni"))
			})

			It("configures the cni_plugin_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniPluginDir).To(Equal("/var/vcap/packages/silk-cni/bin"))
			})
		})

		Context("when External is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "external",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/cni/config/cni"))
			})

			It("configures the cni_plugin_dir", func() {
				manifest, err := product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())
				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())

				Expect(cniPluginDir).To(Equal("/var/vcap/packages/cni/bin"))
			})
		})
	})

	Describe("DNS search domain", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "diego_cell"
			} else {
				instanceGroup = "compute"
			}
		})

		It("configures search_domains on the garden-cni job", func() {
			inputProperties = map[string]interface{}{
				".properties.cf_networking_search_domains": "some-search-domain,another-search-domain",
			}

			manifest, err := product.RenderService.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(Equal([]interface{}{
				"some-search-domain",
				"another-search-domain",
			}))
		})

		It("configures search_domains on the garden-cni job", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(HaveLen(0))
		})
	})

	Describe("Service Discovery For Apps", func() {

		Describe("controller", func() {

			var instanceGroup string

			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "diego_brain"
				} else {
					instanceGroup = "control"
				}
			})

			It("is deployed", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "service-discovery-controller")
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Describe("cell", func() {

			var instanceGroup string

			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "diego_cell"
				} else {
					instanceGroup = "compute"
				}
			})

			It("emits internal routes", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "route_emitter")
				Expect(err).NotTo(HaveOccurred())

				enabled, err := job.Property("internal_routes/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(enabled).To(BeTrue())
			})

		})

	})
})
