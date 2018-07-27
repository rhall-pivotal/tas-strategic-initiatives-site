package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "uaa"
		}
	})

	Describe("route registration", func() {

		It("tags the emitted metrics", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			routeRegistrar, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
			Expect(err).NotTo(HaveOccurred())

			routes, err := routeRegistrar.Property("route_registrar/routes")
			Expect(err).ToNot(HaveOccurred())
			Expect(routes).To(ContainElement(HaveKeyWithValue("tags", map[interface{}]interface{}{
				"component": "uaa",
			})))
		})

	})

	Describe("BPM", func() {
		It("co-locates the BPM job with all diego jobs", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "bpm")
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets bpm.enabled to true", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			manifestJob, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			bpmEnabled, err := manifestJob.Property("bpm/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(bpmEnabled).To(BeTrue())
		})
	})

	Describe("Metrics clients", func() {
		It("apps_metrics has the expected permission scopes", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsScopes, err := uaa.Property("uaa/clients/apps_metrics/scope")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsScopes).To(Equal("cloud_controller.admin,cloud_controller.read,metrics.read,cloud_controller.admin_read_only"))

		})

		It("apps_metrics has the expected redirect uri", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics/redirect-uri")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com,https://metrics.sys.example.com/,https://metrics.sys.example.com/*,https://metrics-previous.sys.example.com,https://metrics-previous.sys.example.com/,https://metrics-previous.sys.example.com/*"))

		})

		It("apps_metrics_processing has the expected permission scopes", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsProcessingScopes, err := uaa.Property("uaa/clients/apps_metrics_processing/scope")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsProcessingScopes).To(Equal("openid,oauth.approvals,doppler.firehose,cloud_controller.admin,cloud_controller.admin_read_only"))

		})

		It("apps_metrics_processing has the expected redirect uri", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics_processing/redirect-uri")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com,https://metrics-previous.sys.example.com"))

		})
	})

})
