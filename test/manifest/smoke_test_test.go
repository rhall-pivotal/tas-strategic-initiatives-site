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
		useExistingOrgProp, err := smokeTest.Property("smoke_tests/use_existing_org")
		Expect(err).NotTo(HaveOccurred())
		Expect(useExistingOrgProp).To(BeTrue())
		org, err := smokeTest.Property("smoke_tests/org")
		Expect(err).NotTo(HaveOccurred())
		Expect(org).To(Equal("system"))

		By("reusing an existing space")
		useExistingSpaceProp, err := smokeTest.Property("smoke_tests/use_existing_space")
		Expect(err).NotTo(HaveOccurred())
		Expect(useExistingSpaceProp).To(BeFalse())

		By("enabling windows tests")
		enableWindowsTestProp, err := smokeTest.Property("smoke_tests/enable_windows_tests")
		Expect(err).NotTo(HaveOccurred())
		Expect(enableWindowsTestProp).To(BeTrue())

		By("setting the windows stack")
		windowsStack, err := smokeTest.Property("smoke_tests/windows_stack")
		Expect(err).NotTo(HaveOccurred())
		Expect(windowsStack).To(Equal("windows"))
	})

	Context("when the operator specifies properties", func() {
		It("configures the smoke-test erand", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.smoke_tests_windows":                       "specified",
				".properties.smoke_tests_windows.specified.org_name":    "banana",
				".properties.smoke_tests_windows.specified.space_name":  "banana",
				".properties.smoke_tests_windows.specified.apps_domain": "banana",
			})
			Expect(err).NotTo(HaveOccurred())

			smokeTest, err := manifest.FindInstanceGroupJob("windows_diego_cell", "smoke_tests_windows")
			Expect(err).NotTo(HaveOccurred())

			By("reusing an existing org")
			useExistingOrgProp, err := smokeTest.Property("smoke_tests/use_existing_org")
			Expect(err).NotTo(HaveOccurred())
			Expect(useExistingOrgProp).To(BeTrue())
			org, err := smokeTest.Property("smoke_tests/org")
			Expect(err).NotTo(HaveOccurred())
			Expect(org).To(Equal("banana"))

			By("reusing an existing space")
			useExistingSpaceProp, err := smokeTest.Property("smoke_tests/use_existing_space")
			Expect(err).NotTo(HaveOccurred())
			Expect(useExistingSpaceProp).To(BeTrue())
			space, err := smokeTest.Property("smoke_tests/space")
			Expect(err).NotTo(HaveOccurred())
			Expect(space).To(Equal("banana"))

			By("enabling windows tests")
			enableWindowsTestProp, err := smokeTest.Property("smoke_tests/enable_windows_tests")
			Expect(err).NotTo(HaveOccurred())
			Expect(enableWindowsTestProp).To(BeTrue())

			By("setting the windows stack")
			windowsStack, err := smokeTest.Property("smoke_tests/windows_stack")
			Expect(err).NotTo(HaveOccurred())
			Expect(windowsStack).To(Equal("windows"))
		})
	})
})
