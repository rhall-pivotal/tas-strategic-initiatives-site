package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("Metrics", func() {
	var (
		getAllInstanceGroups func(planitest.Manifest) []string
	)

	getAllInstanceGroups = func(manifest planitest.Manifest) []string {
		groups, err := manifest.Path("/instance_groups")
		Expect(err).NotTo(HaveOccurred())

		groupList, ok := groups.([]interface{})
		Expect(ok).To(BeTrue())

		names := []string{}
		for _, group := range groupList {
			groupName := group.(map[interface{}]interface{})["name"].(string)

			// ignore VMs that only contain a single placeholder job, i.e. SF-PAS only VMs that are present but non-configurable in PAS build
			jobs, err := manifest.Path(fmt.Sprintf("/instance_groups/name=%s/jobs", groupName))
			Expect(err).NotTo(HaveOccurred())
			if len(jobs.([]interface{})) > 1 {
				names = append(names, groupName)
			}
		}
		Expect(names).NotTo(BeEmpty())
		return names
	}

	Describe("prom scraper", func() {
		It("configures the prom scraper on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

			for _, ig := range instanceGroups {
				scraper, err := manifest.FindInstanceGroupJob(ig, "prom_scraper")
				Expect(err).NotTo(HaveOccurred())

				expectSecureMetrics(scraper)
			}
		})
	})

	Describe("system metric scraper", func() {
		var instanceGroup string
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}
		})

		It("configures the system-metric-scraper", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			metricScraper, err := manifest.FindInstanceGroupJob(instanceGroup, "loggr-system-metric-scraper")
			Expect(err).NotTo(HaveOccurred())

			tlsProps, err := metricScraper.Property("system_metrics/tls")
			Expect(err).ToNot(HaveOccurred())
			Expect(tlsProps).To(HaveKey("ca_cert"))
			Expect(tlsProps).To(HaveKey("cert"))
			Expect(tlsProps).To(HaveKey("key"))

			scrapePort, err := metricScraper.Property("scrape_port")
			Expect(err).ToNot(HaveOccurred())
			Expect(scrapePort).To(Equal(53035))

			expectSecureMetrics(metricScraper)
		})

		It("has a leadership-election job collocated", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			le, err := manifest.FindInstanceGroupJob(instanceGroup, "leadership-election")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := le.Property("port")
			Expect(err).ToNot(HaveOccurred())
			Expect(enabled).To(Equal(7100))
		})
	})

	Describe("metrics discovery", func() {
		It("configures metrics-discovery on all VMs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			instanceGroups := getAllInstanceGroups(manifest)

			for _, ig := range instanceGroups {
				registrar, err := manifest.FindInstanceGroupJob(ig, "metrics-discovery-registrar")
				Expect(err).NotTo(HaveOccurred())

				expectSecureMetrics(registrar)

				agent, err := manifest.FindInstanceGroupJob(ig, "metrics-agent")
				Expect(err).NotTo(HaveOccurred())

				expectSecureMetrics(agent)
			}
		})
	})
})
