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

		var getSecureEndpointWorker = func(manifest planitest.Manifest) planitest.Manifest {
			endpointWorker, err := manifest.FindInstanceGroupJob(workerInstanceGroup, "metric_registrar_secure_endpoint_worker")
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

			It("sets secure endpoint worker defaults", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).ToNot(HaveOccurred())

				secureEndpointWorker := getSecureEndpointWorker(manifest)

				disabled, err := secureEndpointWorker.Property("disabled")
				Expect(err).ToNot(HaveOccurred())
				Expect(disabled).To(Equal(false))

				blacklistedTags, err := secureEndpointWorker.Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("deployment", "job", "index", "id"))

				scrapeInterval, err := secureEndpointWorker.Property("scrape_interval")
				Expect(err).ToNot(HaveOccurred())
				Expect(scrapeInterval).To(Equal("35s"))

				bbsTlsProps, err := secureEndpointWorker.Property("bbs/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(bbsTlsProps).To(HaveKey("ca"))
				Expect(bbsTlsProps).To(HaveKey("cert"))
				Expect(bbsTlsProps).To(HaveKey("key"))

				scrape, err := secureEndpointWorker.Property("scrape")
				Expect(err).ToNot(HaveOccurred())
				Expect(scrape).To(HaveKey("diego_identity_ca"))
				Expect(scrape).To(HaveKey("tls"))

				scrapeTls, err := secureEndpointWorker.Property("scrape/tls")
				Expect(err).ToNot(HaveOccurred())
				Expect(scrapeTls).To(HaveKey("ca"))
				Expect(scrapeTls).To(HaveKey("cert"))
				Expect(scrapeTls).To(HaveKey("key"))
			})

			It("sets up all jobs on different ports", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).ToNot(HaveOccurred())

				allJobs := []planitest.Manifest{
					getEndpointWorker(manifest),
					getSecureEndpointWorker(manifest),
					getOrchestrator(manifest),
					getLogWorker(manifest),
				}

				allPorts := make(map[interface{}]bool)
				for _, job := range allJobs {
					healthCheckPort, err := job.Property("health_check_port")
					Expect(err).ToNot(HaveOccurred())
					_, exists := allPorts[healthCheckPort]
					Expect(exists).To(BeFalse())
					allPorts[healthCheckPort] = true

					orchestrationPort, err := job.Property("orchestration_port")
					Expect(err).ToNot(HaveOccurred())
					_, exists = allPorts[orchestrationPort]
					Expect(exists).To(BeFalse())
					allPorts[orchestrationPort] = true
				}
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

				disabled, err = getSecureEndpointWorker(manifest).Property("disabled")
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

			It("propagates to the secure endpoint worker", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.metric_registrar_scrape_interval_in_seconds": 200,
				})
				Expect(err).ToNot(HaveOccurred())

				scrapeInterval, err := getSecureEndpointWorker(manifest).Property("scrape_interval")
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

				blacklistedTags, err = getSecureEndpointWorker(manifest).Property("blacklisted_tags")
				Expect(err).ToNot(HaveOccurred())
				Expect(blacklistedTags).To(ConsistOf("tag1", "tag2"))
			})
		})
	})
})
