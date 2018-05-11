package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA", func() {
	Describe("Multi-Factor Auth", func() {
		Describe("internal auth", func() {
			Context("when the 'enable multi-factor auth' box is checked", func() {
				It("enables support for multi-factor auth in the uaa manifest", func() {
					manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
						".properties.uaa.internal.enable_mfa_google_authenticator": true,
					})
					Expect(err).NotTo(HaveOccurred())

					agent, err := manifest.FindInstanceGroupJob("uaa", "uaa")
					Expect(err).NotTo(HaveOccurred())

					mfaEnabled, err := agent.Property("login/mfa/enabled")
					Expect(err).NotTo(HaveOccurred())
					Expect(mfaEnabled).To(BeTrue())
				})
			})
			// pending for now because the property is not marked as configurable, TODO: combine with previous test
			XContext("when the 'multi-factor auth issuer' is provided", func() {
				It("renders it on the bosh manifest", func() {
					manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
						".properties.uaa.internal.mfa_google_authenticator_issuer": "some-issuer",
					})
					Expect(err).NotTo(HaveOccurred())

					agent, err := manifest.FindInstanceGroupJob("uaa", "uaa")
					Expect(err).NotTo(HaveOccurred())

					issuer, err := agent.Property("login/providers/googleProvider/config/issuer")
					Expect(err).NotTo(HaveOccurred())
					Expect(issuer).To(Equal("some-issuer"))
				})
			})
		})
		Describe("ldap", func() {
			var ldapProperties map[string]interface{}

			BeforeEach(func() {
				ldapProperties = map[string]interface{}{
					".properties.uaa":                  "ldap",
					".properties.uaa.ldap.url":         "ldaps://some-url",
					".properties.uaa.ldap.credentials": map[string]string{"identity": "some-username", "password": "some-password"},
					".properties.uaa.ldap.search_base": "some-search-base",
				}
			})

			Context("when the 'enable multi-factor auth' box is checked", func() {
				It("enables support for multi-factor auth in the uaa manifest", func() {
					ldapProperties[".properties.uaa.ldap.enable_mfa_google_authenticator"] = true
					manifest, err := product.RenderService.RenderManifest(ldapProperties)
					Expect(err).NotTo(HaveOccurred())

					agent, err := manifest.FindInstanceGroupJob("uaa", "uaa")
					Expect(err).NotTo(HaveOccurred())

					mfaEnabled, err := agent.Property("login/mfa/enabled")
					Expect(err).NotTo(HaveOccurred())
					Expect(mfaEnabled).To(BeTrue())
				})
			})
			// pending for now because the property is not marked as configurable, TODO: combine with previous test
			XContext("when the 'multi-factor auth issuer' is provided", func() {
				It("renders it on the bosh manifest", func() {
					ldapProperties[".properties.uaa.ldap.mfa_google_authenticator_issuer"] = "some-issuer"
					manifest, err := product.RenderService.RenderManifest(ldapProperties)
					Expect(err).NotTo(HaveOccurred())

					agent, err := manifest.FindInstanceGroupJob("uaa", "uaa")
					Expect(err).NotTo(HaveOccurred())

					issuer, err := agent.Property("login/providers/googleProvider/config/issuer")
					Expect(err).NotTo(HaveOccurred())
					Expect(issuer).To(Equal("some-issuer"))
				})
			})
		})
	})
})
