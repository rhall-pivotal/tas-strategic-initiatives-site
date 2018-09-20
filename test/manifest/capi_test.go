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
					Expect(ca).NotTo(BeEmpty())
				}
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

				Expect(stacks).To(ContainElement(map[interface{}]interface{}{
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
			})

		})
	})
})
