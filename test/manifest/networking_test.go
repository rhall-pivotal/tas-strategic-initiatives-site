package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("job colocation", func() {
		It("co-locates bpm", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_diego_cell", "bpm")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("DNS search domain", func() {
		It("configures search_domains on the garden-cni job", func() {
			inputProperties := map[string]interface{}{
				".properties.cf_networking_search_domains": "some-search-domain,another-search-domain",
			}

			manifest, err := product.RenderService.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(Equal([]interface{}{
				"some-search-domain",
				"another-search-domain",
			}))
		})

		It("defaults search_domains to empty", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(HaveLen(0))
		})
	})

	Describe("BOSH DNS Adapter for App Service Discovery", func() {
		It("is colocated with the isolated_diego_cell", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob("isolated_diego_cell", "bosh-dns-adapter")
			Expect(err).NotTo(HaveOccurred())
		})

		It("is turned on by default", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("isolated_diego_cell", "route_emitter")
			Expect(err).NotTo(HaveOccurred())

			enabled, err := job.Property("internal_routes/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(enabled).To(BeTrue())
		})
	})
})
