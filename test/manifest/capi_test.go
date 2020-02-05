package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("CAPI", func() {
	var (
		ccJobs   []Job
		manifest planitest.Manifest
	)

	Describe("common properties", func() {
		BeforeEach(func() {
			if productName == "srt" {
				ccJobs = []Job{
					{
						InstanceGroup: "control",
						Name:          "cloud_controller_ng",
					},
					{
						InstanceGroup: "control",
						Name:          "cloud_controller_worker",
					},
					{
						InstanceGroup: "control",
						Name:          "cloud_controller_clock",
					},
				}
			} else {
				ccJobs = []Job{
					{
						InstanceGroup: "cloud_controller",
						Name:          "cloud_controller_ng",
					},
					{
						InstanceGroup: "cloud_controller_worker",
						Name:          "cloud_controller_worker",
					},
					{
						InstanceGroup: "clock_global",
						Name:          "cloud_controller_clock",
					},
				}
			}
		})

		Context("when the Operator accepts the default values", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets defaults", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					loggingLevel, err := manifestJob.Property("cc/logging_level")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingLevel).To(Equal(string("info")))

					healthCheck, err := manifestJob.Property("cc/default_health_check_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(healthCheck).To(Equal(60))

					diego, err := manifestJob.Property("cc/diego")
					Expect(err).NotTo(HaveOccurred())
					Expect(diego).NotTo(HaveKey("lifecycle_bundles"))

					timeout, err := manifestJob.Property("ccdb/connection_validation_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(timeout).To(Equal(3600))

					timeout, err = manifestJob.Property("ccdb/read_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(timeout).To(Equal(3600))

					address, err := manifestJob.Property("ccdb/address")
					Expect(err).NotTo(HaveOccurred())
					Expect(address).To(Equal("mysql.service.cf.internal"))

					sslVerifyHostname, err := manifestJob.Property("ccdb/ssl_verify_hostname")
					Expect(err).NotTo(HaveOccurred())
					Expect(sslVerifyHostname).To(BeTrue())

					ca, err := manifestJob.Property("ccdb/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(ca).To(BeNil())

					uaaCa, err := manifestJob.Property("uaa/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(uaaCa).NotTo(BeEmpty())

					maxPackageSize, err := manifestJob.Property("cc/packages/max_package_size")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxPackageSize).To(Equal(2147483648))
				}
			})

			It("sets log-cache properties on the cloud_controller_ng job", func() {
				var cloudControllerInstanceGroup string
				if productName == "srt" {
					cloudControllerInstanceGroup = "control"
				} else {
					cloudControllerInstanceGroup = "cloud_controller"
				}

				manifestJob, err := manifest.FindInstanceGroupJob(cloudControllerInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				temporaryUseLogcache, err := manifestJob.Property("cc/temporary_use_logcache")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryUseLogcache).To(Equal(bool(true)))

				Expect(manifestJob.Property("cc/logcache_tls")).Should(HaveKey("certificate"))
				Expect(manifestJob.Property("cc/logcache_tls")).Should(HaveKey("private_key"))
			})

			It("defaults the completed task pruning properties on the cloud_controller_clock job", func() {
				var cloudControllerClockInstanceGroup string
				if productName == "srt" {
					cloudControllerClockInstanceGroup = "control"
				} else {
					cloudControllerClockInstanceGroup = "clock_global"
				}

				manifestJob, err := manifest.FindInstanceGroupJob(cloudControllerClockInstanceGroup, "cloud_controller_clock")
				Expect(err).NotTo(HaveOccurred())

				temporaryUseLogcache, err := manifestJob.Property("cc/completed_tasks/cutoff_age_in_days")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryUseLogcache).To(Equal(31))
			})

			It("does not enable communication to copilot", func() {
				var cloudControllerInstanceGroup string
				if productName == "srt" {
					cloudControllerInstanceGroup = "control"
				} else {
					cloudControllerInstanceGroup = "cloud_controller"
				}

				manifestJob, err := manifest.FindInstanceGroupJob(cloudControllerInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				copilotEnabled, err := manifestJob.Property("copilot/enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(copilotEnabled).To(BeFalse())
			})
		})

		Context("when istio is enabled", func() {
			It("should enable communication to copilot", func() {
				var cloudControllerInstanceGroup string
				if productName == "srt" {
					cloudControllerInstanceGroup = "control"
				} else {
					cloudControllerInstanceGroup = "cloud_controller"
				}

				manifest, err := product.RenderManifest(map[string]interface{}{".properties.istio": "enable"})
				Expect(err).NotTo(HaveOccurred())

				manifestJob, err := manifest.FindInstanceGroupJob(cloudControllerInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				copilotEnabled, err := manifestJob.Property("copilot/enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(copilotEnabled).To(BeTrue())
			})
		})

		Context("when the TLS checkbox is checked", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{".properties.enable_tls_to_internal_pxc": true})
				Expect(err).NotTo(HaveOccurred())
			})

			It("enables TLS to CCDB", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					ca, err := manifestJob.Property("ccdb/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(ca).NotTo(BeEmpty())
				}
			})
		})

		Context("when the Operator accepts the default ZDT deployment updater values", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets defaults", func() {
				var cloudControllerInstanceGroup string
				var clockGlobalInstanceGroup string
				if productName == "srt" {
					cloudControllerInstanceGroup = "control"
					clockGlobalInstanceGroup = "control"
				} else {
					cloudControllerInstanceGroup = "cloud_controller"
					clockGlobalInstanceGroup = "clock_global"
				}

				manifestCloudControllerNgJob, err := manifest.FindInstanceGroupJob(cloudControllerInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				temporaryDisableDeployments, err := manifestCloudControllerNgJob.Property("cc/temporary_disable_deployments")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryDisableDeployments).To(BeFalse())

				manifestCcDeploymentUpdaterJob, err := manifest.FindInstanceGroupJob(clockGlobalInstanceGroup, "cc_deployment_updater")
				Expect(err).NotTo(HaveOccurred())

				temporaryDisableDeployments, err = manifestCcDeploymentUpdaterJob.Property("cc/temporary_disable_deployments")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryDisableDeployments).To(BeFalse())
			})
		})

		Context("when the Operator sets the temporary disable deployments option to true", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".properties.cloud_controller_temporary_disable_deployments": true,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures the subsequent property", func() {
				var cloudControllerInstanceGroup string
				var clockGlobalInstanceGroup string
				if productName == "srt" {
					cloudControllerInstanceGroup = "control"
					clockGlobalInstanceGroup = "control"
				} else {
					cloudControllerInstanceGroup = "cloud_controller"
					clockGlobalInstanceGroup = "clock_global"
				}

				manifestCloudControllerNgJob, err := manifest.FindInstanceGroupJob(cloudControllerInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				temporaryDisableDeployments, err := manifestCloudControllerNgJob.Property("cc/temporary_disable_deployments")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryDisableDeployments).To(BeTrue())

				manifestCcDeploymentUpdaterJob, err := manifest.FindInstanceGroupJob(clockGlobalInstanceGroup, "cc_deployment_updater")
				Expect(err).NotTo(HaveOccurred())

				temporaryDisableDeployments, err = manifestCcDeploymentUpdaterJob.Property("cc/temporary_disable_deployments")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryDisableDeployments).To(BeTrue())
			})
		})

		Context("when the Operator sets the completed task cutoff age in days to custom values", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".properties.cloud_controller_completed_tasks_cutoff_age_in_days": 32,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures the subsequent property", func() {
				var cloudControllerClockInstanceGroup string
				if productName == "srt" {
					cloudControllerClockInstanceGroup = "control"
				} else {
					cloudControllerClockInstanceGroup = "clock_global"
				}

				manifestJob, err := manifest.FindInstanceGroupJob(cloudControllerClockInstanceGroup, "cloud_controller_clock")
				Expect(err).NotTo(HaveOccurred())

				temporaryUseLogcache, err := manifestJob.Property("cc/completed_tasks/cutoff_age_in_days")
				Expect(err).NotTo(HaveOccurred())
				Expect(temporaryUseLogcache).To(Equal(32))
			})
		})

		Context("when the Operator sets CC logging level to debug", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".properties.cc_logging_level": "debug",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures logging level to debug", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					loggingLevel, err := manifestJob.Property("cc/logging_level")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingLevel).To(Equal(string("debug")))
				}
			})
		})

		Context("when the Operator sets the Database Connection Validation Timeout", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".properties.ccdb_connection_validation_timeout": 200,
					".properties.ccdb_read_timeout":                  200,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures the timeouts on the ccdb", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					timeout, err := manifestJob.Property("ccdb/connection_validation_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(timeout).To(Equal(200))

					timeout, err = manifestJob.Property("ccdb/read_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(timeout).To(Equal(200))
				}
			})
		})

		Context("when the Operator sets the Default Health Check Timeout", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".properties.cloud_controller_default_health_check_timeout": 120,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("passes the value to CC jobs", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					healthCheck, err := manifestJob.Property("cc/default_health_check_timeout")
					Expect(err).NotTo(HaveOccurred())

					Expect(healthCheck).To(Equal(120))
				}
			})
		})

		Context("when the Operator sets an Insecure Registry list", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".diego_cell.insecure_docker_registry_list": "item1,item2,item3",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("passes the value to CC jobs", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					insecureDockerRegistryList, err := manifestJob.Property("cc/diego/insecure_docker_registry_list")
					Expect(err).NotTo(HaveOccurred())

					Expect(insecureDockerRegistryList).To(Equal([]interface{}{"item1", "item2", "item3"}))
				}
			})
		})

		Context("when the Operator sets a staging timeout", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".cloud_controller.staging_timeout_in_seconds": 1000,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("passes the value to CC jobs", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					insecureDockerRegistryList, err := manifestJob.Property("cc/staging_timeout_in_seconds")
					Expect(err).NotTo(HaveOccurred())

					Expect(insecureDockerRegistryList).To(Equal(1000))
				}
			})
		})

		Context("when the Operator sets a max package size", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".cloud_controller.max_package_size": 5368709120,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("passes the value to CC jobs", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					maxPackageSize, err := manifestJob.Property("cc/packages/max_package_size")
					Expect(err).NotTo(HaveOccurred())

					Expect(maxPackageSize).To(Equal(5368709120))
				}
			})
		})

		Context("when the Operator sets a max disk quota for an app", func() {
			BeforeEach(func() {
				var err error
				manifest, err = product.RenderManifest(map[string]interface{}{
					".cloud_controller.max_disk_quota_app": 10240,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("passes the value to CC jobs", func() {
				for _, job := range ccJobs {
					manifestJob, err := manifest.FindInstanceGroupJob(job.InstanceGroup, job.Name)
					Expect(err).NotTo(HaveOccurred())

					maxPackageSize, err := manifestJob.Property("cc/maximum_app_disk_in_mb")
					Expect(err).NotTo(HaveOccurred())

					Expect(maxPackageSize).To(Equal(10240))
				}
			})
		})
	})

	Describe("api", func() {

		var instanceGroup string

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "cloud_controller"
			}

			var err error
			manifest, err = product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("keeps the docs link up-to-date", func() {
			api, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
			Expect(err).NotTo(HaveOccurred())

			description, err := api.Property("description")
			Expect(err).NotTo(HaveOccurred())
			Expect(description).To(MatchRegexp(`https://docs.pivotal.io/pivotalcf/\d+-\d+/pcf-release-notes/runtime-rn.html`))
		})

		It("sets defaults on the udp forwarder", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			udpForwarder, err := manifest.FindInstanceGroupJob(instanceGroup, "loggr-udp-forwarder")
			Expect(err).NotTo(HaveOccurred())

			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("ca"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("cert"))
			Expect(udpForwarder.Property("loggregator/tls")).Should(HaveKey("key"))
		})

		It("uses cflinuxfs3 for the docker staging stack", func() {
			api, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
			Expect(err).NotTo(HaveOccurred())

			description, err := api.Property("cc/diego/docker_staging_stack")
			Expect(err).NotTo(HaveOccurred())
			Expect(description).To(Equal("cflinuxfs3"))
		})


		Describe("tls routing", func() {
			It("configures the route registrar to use tls", func() {
				routeRegistrarJob, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
				Expect(err).NotTo(HaveOccurred())

				tlsPort, err := routeRegistrarJob.Property("route_registrar/routes/name=api/tls_port")
				Expect(err).NotTo(HaveOccurred())
				Expect(tlsPort).To(Equal(9024))

				certAltName, err := routeRegistrarJob.Property("route_registrar/routes/name=api/server_cert_domain_san")
				Expect(err).NotTo(HaveOccurred())
				Expect(certAltName).To(Equal("cloud-controller-ng.service.cf.internal"))
			})

			It("configures the cloud controller tls certs", func() {
				cloudControllerJob, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				Expect(cloudControllerJob.Property("cc/public_tls")).Should(HaveKey("ca_cert"))
				Expect(cloudControllerJob.Property("cc/public_tls")).Should(HaveKey("certificate"))
				Expect(cloudControllerJob.Property("cc/public_tls")).Should(HaveKey("private_key"))
			})
		})

		Describe("stacks", func() {

			It("defines stacks", func() {
				cc, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				stacks, err := cc.Property("cc/stacks")
				Expect(err).NotTo(HaveOccurred())

				Expect(stacks).NotTo(ContainElement(map[interface{}]interface{}{
					"name":        "cflinuxfs2",
					"description": "Cloud Foundry Linux-based filesystem - Ubuntu Trusty 14.04 LTS",
				}))
				Expect(stacks).To(ContainElement(map[interface{}]interface{}{
					"name":        "cflinuxfs3",
					"description": "Cloud Foundry Linux-based filesystem - Ubuntu Bionic 18.04 LTS",
				}))
				Expect(stacks).To(ContainElement(map[interface{}]interface{}{
					"name":        "windows2012R2",
					"description": "Microsoft Windows / .Net 64 bit",
				}))
				Expect(stacks).To(ContainElement(map[interface{}]interface{}{
					"name":        "windows2016",
					"description": "Microsoft Windows 2016",
				}))
				Expect(stacks).To(ContainElement(map[interface{}]interface{}{
					"name":        "windows",
					"description": "Windows Server",
				}))

				defaultStack, err := cc.Property("cc/default_stack")
				Expect(err).NotTo(HaveOccurred())
				Expect(defaultStack).To(Equal("cflinuxfs3"))
			})
		})

		Describe("Database Encryption Keys", func() {
			It("sets the encryption keys in cloud controller job", func() {
				cloudControllerJob, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				databaseEncryptionKeys, err := cloudControllerJob.Property("cc/database_encryption/keys")
				Expect(err).NotTo(HaveOccurred())
				Expect(databaseEncryptionKeys).To(Equal([]interface{}{}))
			})

			Context("when the encryption keys are provided", func() {
				It("sets the encryption keys in cloud controller job", func() {
					manifest, err := product.RenderManifest((map[string]interface{}{
						".properties.cloud_controller_encryption_keys": []map[string]interface{}{
							{
								"encryption_key": map[string]interface{}{
									"secret": "some-encryption-key",
								},
								"label":   "some internal key display name",
								"primary": true,
							},
							{
								"encryption_key": map[string]interface{}{
									"secret": "old-encryption-key",
								},
								"label":   "old internal key display name",
								"primary": false,
							},
						},
					}))
					Expect(err).NotTo(HaveOccurred())
					cloudControllerJob, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					databaseEncryptionKeys, err := cloudControllerJob.Property("cc/database_encryption/keys")
					Expect(err).NotTo(HaveOccurred())
					Expect(databaseEncryptionKeys).To(HaveLen(2))
				})
			})
		})
	})
})
