package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	It("has a push-usage-service-client used to push the app and for bbr", func() {
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
			pushUsageService, err := manifest.FindInstanceGroupJob(pushUsageServiceInstanceGroup, "push-usage-service")

			Expect(err).NotTo(HaveOccurred())

			_, err = pushUsageService.Property("cf/admin_username")
			Expect(err).To(HaveOccurred())

			_, err = pushUsageService.Property("cf/admin_password")
			Expect(err).To(HaveOccurred())

			clientID, err := pushUsageService.Property("cf/push_client_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(clientID).To(Equal(pushUsageServiceClientID))

			_, err = pushUsageService.Property("cf/push_client_secret")
			Expect(err).NotTo(HaveOccurred())
		})

		By("configuring the bbr-usage-servicedb", func() {
			bbrUsageServiceDB, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-usage-servicedb")
			Expect(err).NotTo(HaveOccurred())

			_, err = bbrUsageServiceDB.Property("cf/admin_username")
			Expect(err).To(HaveOccurred())

			_, err = bbrUsageServiceDB.Property("cf/admin_password")
			Expect(err).To(HaveOccurred())

			clientID, err := bbrUsageServiceDB.Property("cf/client_id")
			Expect(err).NotTo(HaveOccurred())
			Expect(clientID).To(Equal(pushUsageServiceClientID))

			_, err = bbrUsageServiceDB.Property("cf/client_secret")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
