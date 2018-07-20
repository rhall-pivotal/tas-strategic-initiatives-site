package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CAPI", func() {
	var instanceGroup string

	Context("insecure docker registry list", func() {
		Context("for the cloud_controller_ng job", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "cloud_controller"
				}
			})

			It("configures a insecure_docker_registry_list", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				cloudControllerClock, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				insecureDockerRegistryList, err := cloudControllerClock.Property("cc/diego/insecure_docker_registry_list")
				Expect(err).NotTo(HaveOccurred())

				Expect(insecureDockerRegistryList).To(Equal([]interface{}{"item1", "item2", "item3"}))
			})
		})

		Context("for the cloud_controller_worker job", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "cloud_controller_worker"
				}
			})

			It("configures a insecure_docker_registry_list", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				cloudControllerClock, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_worker")
				Expect(err).NotTo(HaveOccurred())

				insecureDockerRegistryList, err := cloudControllerClock.Property("cc/diego/insecure_docker_registry_list")
				Expect(err).NotTo(HaveOccurred())

				Expect(insecureDockerRegistryList).To(Equal([]interface{}{"item1", "item2", "item3"}))
			})
		})

		Context("for the cloud_controller_clock job", func() {
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "clock_global"
				}
			})

			It("configures a insecure_docker_registry_list", func() {
				manifest, err := product.RenderService.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				cloudControllerClock, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_clock")
				Expect(err).NotTo(HaveOccurred())

				insecureDockerRegistryList, err := cloudControllerClock.Property("cc/diego/insecure_docker_registry_list")
				Expect(err).NotTo(HaveOccurred())

				Expect(insecureDockerRegistryList).To(Equal([]interface{}{"item1", "item2", "item3"}))
			})
		})
	})
})
