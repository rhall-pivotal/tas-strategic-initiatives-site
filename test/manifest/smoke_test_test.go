package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SmokeTest", func() {
	It("contributes a smoke test instance group with the errand", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		smokeTest, err := manifest.FindInstanceGroupJob("windows_diego_cell", "smoke_tests_windows")
		Expect(err).NotTo(HaveOccurred())

		By("reusing an existing org")
		org, err := smokeTest.Property("smoke_tests/organization")
		Expect(err).NotTo(HaveOccurred())
		Expect(org).To(Equal("system"))

		By("setting the windows stack")
		windowsStack, err := smokeTest.Property("smoke_tests/windows_stack")
		Expect(err).NotTo(HaveOccurred())
		Expect(windowsStack).To(Equal("windows"))
	})

	Context("when the operator specifies properties", func() {
		It("configures the smoke-test erand", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.smoke_tests_windows":                      "specified",
				".properties.smoke_tests_windows.specified.org_name":   "banana",
				".properties.smoke_tests_windows.specified.space_name": "banana",
			})
			Expect(err).NotTo(HaveOccurred())

			smokeTest, err := manifest.FindInstanceGroupJob("windows_diego_cell", "smoke_tests_windows")
			Expect(err).NotTo(HaveOccurred())

			By("reusing an existing org")
			org, err := smokeTest.Property("smoke_tests/organization")
			Expect(err).NotTo(HaveOccurred())
			Expect(org).To(Equal("banana"))

			By("reusing an existing space")
			space, err := smokeTest.Property("smoke_tests/space")
			Expect(err).NotTo(HaveOccurred())
			Expect(space).To(Equal("banana"))

			By("setting the windows stack")
			windowsStack, err := smokeTest.Property("smoke_tests/windows_stack")
			Expect(err).NotTo(HaveOccurred())
			Expect(windowsStack).To(Equal("windows"))
		})
	})
})
