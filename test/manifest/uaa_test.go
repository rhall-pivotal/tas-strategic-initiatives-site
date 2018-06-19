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

})
