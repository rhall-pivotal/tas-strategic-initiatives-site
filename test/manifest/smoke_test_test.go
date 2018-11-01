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
	})
})
