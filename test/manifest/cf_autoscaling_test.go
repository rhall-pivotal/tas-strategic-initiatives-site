package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CF Autoscaling", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
	})

	It("sets the organization and space for the test-autoscale errand", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "test-autoscaling")
		Expect(err).NotTo(HaveOccurred())

		space, err := job.Property("autoscale/space")
		Expect(err).ToNot(HaveOccurred())
		Expect(space).To(Equal("autoscaling"))

		org, err := job.Property("autoscale/organization")
		Expect(err).ToNot(HaveOccurred())
		Expect(org).To(Equal("system"))
	})

	It("keeps the doc link up-to-date", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
		Expect(err).NotTo(HaveOccurred())

		url, err := job.Property("autoscale/marketplace_documentation_url")
		Expect(err).ToNot(HaveOccurred())
		Expect(url).To(MatchRegexp(`https://docs.pivotal.io/pivotalcf/\d+-\d+/appsman-services/autoscaler/using-autoscaler.html`))
	})

	It("enables notifications by default", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
		Expect(err).NotTo(HaveOccurred())

		property, err := job.Property("autoscale/enable_notifications")
		Expect(err).ToNot(HaveOccurred())
		Expect(property).To(BeTrue())
	})

	It("sets the log_level to error", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
		Expect(err).NotTo(HaveOccurred())

		logLevel, err := job.Property("autoscale/log_level")
		Expect(err).ToNot(HaveOccurred())
		Expect(logLevel).To(Equal("error"))
	})

	Context("when the user enables verbose logging", func() {
		It("enables verbose logging", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.autoscale_enable_verbose_logging": true,
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
			Expect(err).NotTo(HaveOccurred())

			property, err := job.Property("autoscale/enable_verbose_logging")
			Expect(err).ToNot(HaveOccurred())
			Expect(property).To(BeTrue())
		})
		It("sets the log_level to info", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.autoscale_enable_verbose_logging": true,
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
			Expect(err).NotTo(HaveOccurred())

			logLevel, err := job.Property("autoscale/log_level")
			Expect(err).ToNot(HaveOccurred())
			Expect(logLevel).To(Equal("info"))
		})
	})

	Context("when the user disables connection pooling", func() {
		It("sets the autoscale api to disable connection pooling", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.autoscale_api_disable_connection_pooling": true,
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
			Expect(err).NotTo(HaveOccurred())

			property, err := job.Property("autoscale/api/disable_connection_pooling")
			Expect(err).ToNot(HaveOccurred())
			Expect(property).To(BeTrue())
		})
	})

	Context("when the user disables notifications", func() {
		It("sets the autoscale api to disable connection pooling", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.autoscale_enable_notifications": false,
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-autoscaler")
			Expect(err).NotTo(HaveOccurred())

			property, err := job.Property("autoscale/enable_notifications")
			Expect(err).ToNot(HaveOccurred())
			Expect(property).To(BeFalse())
		})
	})

	Describe("Backup and Restore", func() {
		Context("on the backup_restore instance group", func() {
			It("templates the deploy-autoscaler job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob("backup_restore", "deploy-autoscaler")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
