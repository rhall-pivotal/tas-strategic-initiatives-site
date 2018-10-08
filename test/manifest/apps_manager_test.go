package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apps Manager", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
	})

	Describe("BPM", func() {
		var appsManagerJobs []Job

		BeforeEach(func() {
			if productName == "srt" {
				appsManagerJobs = []Job{
					{
						InstanceGroup: "control",
						Name:          "push-apps-manager",
					},
				}
			} else {
				appsManagerJobs = []Job{
					{
						InstanceGroup: "clock_global",
						Name:          "push-apps-manager",
					},
				}
			}
		})

		It("co-locates the BPM job with all apps manager jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, appsManagerJob := range appsManagerJobs {
				_, err = manifest.FindInstanceGroupJob(appsManagerJob.InstanceGroup, "bpm")
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("sets bpm.enabled to true", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, appsManagerJob := range appsManagerJobs {
				manifestJob, err := manifest.FindInstanceGroupJob(appsManagerJob.InstanceGroup, appsManagerJob.Name)
				Expect(err).NotTo(HaveOccurred())

				bpmEnabled, err := manifestJob.Property("bpm/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(bpmEnabled).To(BeTrue())
			}
		})
	})

	Describe("Memory", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
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
				manifest, err := product.RenderManifest(map[string]interface{}{
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

	Describe("Polling intervals", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			pollInterval, err := appsManager.Property("apps_manager/poll_interval")
			Expect(err).NotTo(HaveOccurred())
			Expect(pollInterval).To(Equal(30))

			appPollInterval, err := appsManager.Property("apps_manager/app_poll_interval")
			Expect(err).NotTo(HaveOccurred())
			Expect(appPollInterval).To(Equal(10))
		})

		Context("when the operator specifies a poll interval", func() {
			It("applies them", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_app_poll_interval": 333,
					".properties.push_apps_manager_poll_interval":     666,
				})
				Expect(err).NotTo(HaveOccurred())

				appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				pollInterval, err := appsManager.Property("apps_manager/poll_interval")
				Expect(err).NotTo(HaveOccurred())
				Expect(pollInterval).To(Equal(666))

				appPollInterval, err := appsManager.Property("apps_manager/app_poll_interval")
				Expect(err).NotTo(HaveOccurred())
				Expect(appPollInterval).To(Equal(333))
			})
		})
	})
})
