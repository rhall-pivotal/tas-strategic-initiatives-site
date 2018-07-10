package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	var instanceGroup string

	Describe("notifications", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}
		})

		It("has a notifications job with default CF notification template", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			notifications, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-notifications")
			Expect(err).NotTo(HaveOccurred())

			template, err := notifications.Property("notifications/default_template")
			Expect(err).ToNot(HaveOccurred())
			Expect(template).To(ContainSubstring("CF Notification: {{.Subject}}"))
		})
	})
})
