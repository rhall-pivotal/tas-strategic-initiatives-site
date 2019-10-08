package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Usage Service", func() {
	var (
		uaaInstanceGroup              string
		pushUsageServiceInstanceGroup string
	)

	BeforeEach(func() {
		if productName == "srt" {
			uaaInstanceGroup = "control"
			pushUsageServiceInstanceGroup = "control"

		} else {
			uaaInstanceGroup = "uaa"
			pushUsageServiceInstanceGroup = "clock_global"
		}
	})

	Describe("Deploying Usage Service", func() {
		It("has a push-usage-service-client used to push the app", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			pushUsageServiceClientID := "push_usage_service"

			By("creating a uaa client", func() {
				uaa, err := manifest.FindInstanceGroupJob(uaaInstanceGroup, "uaa")
				Expect(err).NotTo(HaveOccurred())

				client, err := uaa.Property("uaa/clients/" + pushUsageServiceClientID)
				Expect(err).NotTo(HaveOccurred())

				Expect(client).To(HaveKeyWithValue("authorities", "cloud_controller.admin"))
				Expect(client).To(HaveKeyWithValue("authorized-grant-types", "client_credentials"))

				_, err = uaa.Property("uaa/clients/push_usage_service/secret")
				Expect(err).NotTo(HaveOccurred())
			})

			By("configuring the push-usage-service", func() {
				job, err := manifest.FindInstanceGroupJob(pushUsageServiceInstanceGroup, "push-usage-service")
				Expect(err).NotTo(HaveOccurred())

				testPushUsageServiceProperties(job)
				Expect(job.Path("/provides/app-usage-internal")).To(Equal(map[interface{}]interface{}{"as": "app-usage-internal"}))
			})
		})

		Describe("Cutoff age in days", func() {
			It("uses the spec defaults", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(pushUsageServiceInstanceGroup, "push-usage-service")
				Expect(err).NotTo(HaveOccurred())

				Expect(job.Property("app_usage_service/cutoff_age_in_days")).To(Equal(365))
			})

			Context("when the operator specifies cutoff value manually", func() {
				It("applies them", func() {
					manifest, err := product.RenderManifest(map[string]interface{}{
						".properties.push_usage_service_cutoff_age_in_days": 100,
					})
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(pushUsageServiceInstanceGroup, "push-usage-service")
					Expect(err).NotTo(HaveOccurred())

					Expect(job.Property("app_usage_service/cutoff_age_in_days")).To(Equal(100))
				})
			})
		})
	})

	Describe("Backup and Restore", func() {
		Context("on the backup_restore instance group", func() {
			It("configures the push-usage-service job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "push-usage-service")
				Expect(err).NotTo(HaveOccurred())

				testPushUsageServiceProperties(job)
				Expect(job.Path("/provides/app-usage-internal")).To(Equal(map[interface{}]interface{}{"as": "ignore-me"}))
			})

			It("configures the bbr-usage-servicedb job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				bbrUsageService, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-usage-servicedb")
				Expect(err).NotTo(HaveOccurred())

				Expect(bbrUsageService.Path("/consumes/app-usage-internal")).To(Equal(map[interface{}]interface{}{"from": "app-usage-internal"}))
			})
		})
	})
})

func testPushUsageServiceProperties(pushUsageService planitest.Manifest) {
	_, err := pushUsageService.Property("cf/admin_username")
	Expect(err).To(HaveOccurred())

	_, err = pushUsageService.Property("cf/admin_password")
	Expect(err).To(HaveOccurred())

	Expect(pushUsageService.Property("cf/push_client_id")).To(Equal("push_usage_service"))

	_, err = pushUsageService.Property("cf/push_client_secret")
	Expect(err).NotTo(HaveOccurred())
}
