package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("Container networking", func() {
		var (
			inputProperties         map[string]interface{}
			cellInstanceGroup       string
			controllerInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "ert" {
				cellInstanceGroup = "diego_cell"
				controllerInstanceGroup = "diego_database"
			} else {
				cellInstanceGroup = "compute"
				controllerInstanceGroup = "control"
			}
		})

		Describe("policy server", func() {
			It("uses the correct database host", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				host, err := job.Property("database/host")
				Expect(err).NotTo(HaveOccurred())
				Expect(host).To(Equal("mysql.service.cf.internal"))

				databaseLink, err := job.Path("/consumes/database")
				Expect(err).NotTo(HaveOccurred())
				Expect(databaseLink).To(Equal("nil"))
			})
		})

		Context("when the operator configures database connection timeout for CNI plugin", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.cf_networking_database_connection_timeout": 250,
				}
			})

			It("sets the manifest database connection timeout properties for the cf networking jobs to be 250", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				policyServerJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				policyServerConnectTimeoutSeconds, err := policyServerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerConnectTimeoutSeconds).To(Equal(250))

				policyServerInternalJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server-internal")
				Expect(err).NotTo(HaveOccurred())

				policyServerInternalConnectTimeoutSeconds, err := policyServerInternalJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerInternalConnectTimeoutSeconds).To(Equal(250))

				silkControllerJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				silkControllerConnectTimeoutSeconds, err := silkControllerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(silkControllerConnectTimeoutSeconds).To(Equal(250))
			})
		})

		Context("when Silk is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{}
			})

			It("configures the cni_config_dir and cni_plugin_dir", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/silk-cni/config/cni"))

				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniPluginDir).To(Equal("/var/vcap/packages/silk-cni/bin"))
			})

			It("configures TLS to the internal database", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				tlsEnabled, err := job.Property("database/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(tlsEnabled).To(BeTrue())

				caCert, err := job.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).NotTo(BeEmpty())
			})

			It("uses the correct database host", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				host, err := job.Property("database/host")
				Expect(err).NotTo(HaveOccurred())
				Expect(host).To(Equal("mysql.service.cf.internal"))

				databaseLink, err := job.Path("/consumes/database")
				Expect(err).NotTo(HaveOccurred())
				Expect(databaseLink).To(Equal("nil"))
			})

		})

		Context("when External is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "external",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/cni/config/cni"))

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

			manifest, err := product.RenderManifest(inputProperties)
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
			manifest, err := product.RenderManifest(nil)
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
				manifest, err := product.RenderManifest(nil)
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

			It("co-locates the bosh-dns-adapter and bpm", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "bpm")
				Expect(err).NotTo(HaveOccurred())
			})

			It("emits internal routes", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "route_emitter")
				Expect(err).NotTo(HaveOccurred())

				enabled, err := job.Property("internal_routes/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(enabled).To(BeTrue())
			})

			Context("when internal domain is empty", func() {
				It("defaults internal domain to apps.internal", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("internal_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(ConsistOf("apps.internal."))
				})
			})

			Context("when internal domain is configured", func() {
				var (
					inputProperties map[string]interface{}
				)

				It("sets internal domains to the provided internal domains", func() {
					inputProperties = map[string]interface{}{
						".properties.cf_networking_internal_domain": "some-internal-domain",
					}
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("internal_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"some-internal-domain",
					}))
				})
			})
		})

		Describe("api", func() {
			var instanceGroup string
			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "cloud_controller"
				} else {
					instanceGroup = "control"
				}
			})

			Context("when internal domain is empty", func() {
				It("adds apps.internal to app domains", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("app_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"apps.example.com",
						map[interface{}]interface{}{
							"name":     "apps.internal.",
							"internal": true,
						},
					}))
				})
			})

			Context("when internal domain is configured", func() {
				var (
					inputProperties map[string]interface{}
				)

				It("adds internal domains to app domains", func() {
					inputProperties = map[string]interface{}{
						".properties.cf_networking_internal_domain": "some-internal-domain",
					}
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("app_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"apps.example.com",
						map[interface{}]interface{}{
							"name":     "some-internal-domain",
							"internal": true,
						},
					}))
				})
			})
		})
	})
})
