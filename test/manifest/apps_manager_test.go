package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apps Manager", func() {

	Describe("Memory", func() {

		var instanceGroup string

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}
		})

		It("uses the spec defaults", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			appsManagerMemory, err := appsManager.Property("apps_manager/memory")
			Expect(err).NotTo(HaveOccurred())
			Expect(appsManagerMemory).To(BeNil())

			invitationsMemory, err := appsManager.Property("invitations/memory")
			Expect(err).NotTo(HaveOccurred())
			Expect(invitationsMemory).To(BeNil())
		})

		Context("when the operator specifies memory limits", func() {

			It("applies them", func() {
				manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_memory":             1024,
					".properties.push_apps_manager_invitations_memory": 2048,
				})
				Expect(err).NotTo(HaveOccurred())

				appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				appsManagerMemory, err := appsManager.Property("apps_manager/memory")
				Expect(err).NotTo(HaveOccurred())
				Expect(appsManagerMemory).To(Equal(1024))

				invitationsMemory, err := appsManager.Property("invitations/memory")
				Expect(err).NotTo(HaveOccurred())
				Expect(invitationsMemory).To(Equal(2048))
			})

		})

	})

})
