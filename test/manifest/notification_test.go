package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Notifications", func() {
	var (
		instanceGroup   string
		inputProperties map[string]interface{}
	)

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
		inputProperties = map[string]interface{}{}
	})

	It("has a notifications job with default properties", func() {
		manifest, err := product.RenderManifest(inputProperties)
		Expect(err).NotTo(HaveOccurred())

		notifications, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-notifications")
		Expect(err).NotTo(HaveOccurred())

		template, err := notifications.Property("notifications/default_template")
		Expect(err).ToNot(HaveOccurred())
		Expect(template).To(ContainSubstring("CF Notification: {{.Subject}}"))

		caCert, err := notifications.Property("notifications/database/ca_cert")
		Expect(err).ToNot(HaveOccurred())
		Expect(caCert).To(BeNil())
	})

	Context("when the TLS checkbox is checked", func() {
		BeforeEach(func() {
			inputProperties = map[string]interface{}{
				".properties.enable_tls_to_internal_pxc": true,
			}
		})

		It("enables TLS", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			notifications, err := manifest.FindInstanceGroupJob(instanceGroup, "deploy-notifications")
			Expect(err).NotTo(HaveOccurred())

			caCert, err := notifications.Property("notifications/database/ca_cert")
			Expect(err).ToNot(HaveOccurred())
			Expect(caCert).NotTo(BeEmpty())
		})
	})
})
