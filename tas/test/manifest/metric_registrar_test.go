package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Metric Registrar", func() {
	Describe("metric registrar", func() {
		var (
			workerInstanceGroup       string
			orchestratorInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "srt" {
				workerInstanceGroup = "control"
				orchestratorInstanceGroup = "control"
			} else {
				workerInstanceGroup = "doppler"
				orchestratorInstanceGroup = "clock_global"
			}
		})

		var getOrchestrator = func(manifest planitest.Manifest) planitest.Manifest {
			orchestrator, err := manifest.FindInstanceGroupJob(orchestratorInstanceGroup, "metric_registrar_orchestrator")
			Expect(err).ToNot(HaveOccurred())
			return orchestrator
		}

		var getLogWorker = func(manifest planitest.Manifest) planitest.Manifest {
			logWorker, err := manifest.FindInstanceGroupJob(workerInstanceGroup, "metric_registrar_log_worker")
			Expect(err).ToNot(HaveOccurred())
			return logWorker
		}

		var getEndpointWorker = func(manifest planitest.Manifest) planitest.Manifest {
			endpointWorker, err := manifest.FindInstanceGroupJob(workerInstanceGroup, "metric_registrar_endpoint_worker")
			Expect(err).ToNot(HaveOccurred())
			return endpointWorker
		}

		Context("defaults", func() {
			It("sets orchestrator defaults", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).ToNot(HaveOccurred())

				disabled, err := getOrchestrator(manifest).Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(false))
			})

			It("sets log worker defaults", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).ToNot(HaveOccurred())
				logWorker := getLogWorker(manifest)

				disabled, err := logWorker.Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(false))

				blacklistedTags, err := logWorker.Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("deployment", "job", "index", "id"))
			})

			It("sets endpoint worker defaults", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).ToNot(HaveOccurred())

				endpointWorker := getEndpointWorker(manifest)

				disabled, err := endpointWorker.Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(false))

				blacklistedTags, err := endpointWorker.Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("deployment", "job", "index", "id"))

				scrapeInterval, err := endpointWorker.Property("scrape_interval")
				Expect(err).ToNot(HaveOccurred())
				Expect(scrapeInterval).To(Equal("35s"))
			})
		})

		When("operator disables metric registrar", func() {
			It("disables the orchestrator and workers", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_enabled": false,
				})
				Expect(err).ToNot(HaveOccurred())

				disabled, err := getOrchestrator(manifest).Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(true))

				disabled, err = getLogWorker(manifest).Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(true))

				disabled, err = getEndpointWorker(manifest).Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(true))
			})
		})

		When("operator sets scrape interval", func() {
			It("propagates to the endpoint worker", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_scrape_interval_in_seconds": 200,
				})
				Expect(err).ToNot(HaveOccurred())

				scrapeInterval, err := getEndpointWorker(manifest).Property("scrape_interval")
				Expect(err).ToNot(HaveOccurred())
				Expect(scrapeInterval).To(Equal("200s"))
			})

			It("errors if set below minimum", func() {
				_, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_scrape_interval_in_seconds": 14,
				})
				Expect(err).To(HaveOccurred())
			})

			It("errors if set above maximum", func() {
				_, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_scrape_interval_in_seconds": 601,
				})
				Expect(err).To(HaveOccurred())
			})
		})

		When("operator sets blacklisted tags", func() {
			It("propagates to the workers", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_blacklisted_tags": "tag1,tag2",
				})
				Expect(err).ToNot(HaveOccurred())

				blacklistedTags, err := getLogWorker(manifest).Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("tag1", "tag2"))

				blacklistedTags, err = getEndpointWorker(manifest).Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("tag1", "tag2"))
			})
		})
	})
})
