package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CredHub", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "credhub"
		}
	})

	Describe("encryption keys", func() {

		Context("when there is a single internal key", func() {

			It("configures credhub with the key", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				keys, err := credhub.Property("credhub/encryption/keys")
				Expect(err).ToNot(HaveOccurred())
				Expect(keys).To(HaveLen(1))

				key := keys.([]interface{})[0].(map[interface{}]interface{})
				Expect(key["provider_name"]).To(Equal("internal-provider"))
				Expect(key["key_properties"]).To(HaveKeyWithValue("encryption_password", ContainSubstring("credhub_key_encryption_passwords/0/key.value")))
				Expect(key["active"]).To(BeTrue())
			})

		})

		Context("when there is an additional HSM key set as primary", func() {

			It("configures credhub with the keys, with the HSM key marked as active", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.credhub_key_encryption_passwords": []map[string]interface{}{
						{
							"key": map[string]interface{}{
								"secret": "some-credhub-password",
							},
							"name":     "internal key display name",
							"primary":  false,
							"provider": "internal",
						},
						{
							"key": map[string]interface{}{
								"secret": "hsm-provider-key-name",
							},
							"name":     "hsm key display name",
							"primary":  true,
							"provider": "hsm",
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				keys, err := credhub.Property("credhub/encryption/keys")
				Expect(err).ToNot(HaveOccurred())
				Expect(keys).To(HaveLen(2))

				internalKey := keys.([]interface{})[0].(map[interface{}]interface{})
				Expect(internalKey["provider_name"]).To(Equal("internal-provider"))
				Expect(internalKey["key_properties"]).To(HaveKeyWithValue("encryption_password", ContainSubstring("credhub_key_encryption_passwords/0/key.value")))
				Expect(internalKey["active"]).To(BeFalse())

				hsmKey := keys.([]interface{})[1].(map[interface{}]interface{})
				Expect(hsmKey["provider_name"]).To(Equal("hsm-provider"))
				Expect(hsmKey["key_properties"]).To(HaveKeyWithValue("encryption_key_name", ContainSubstring("credhub_key_encryption_passwords/1/key.value")))
				Expect(hsmKey["active"]).To(BeTrue())
			})

		})

	})

	Describe("permissions", func() {

		It("grants permissions to the credhub-service-broker tile", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			permissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).ToNot(HaveOccurred())

			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/credhub-clients/*",
				"actors":     []interface{}{"uaa-client:credhub-service-broker"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))

			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/credhub-service-broker/*",
				"actors":     []interface{}{"uaa-client:credhub-service-broker"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))
		})

		It("grants permissions to the services_credhub_client", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			permissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).ToNot(HaveOccurred())

			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/c/*",
				"actors":     []interface{}{"uaa-client:services_credhub_client"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))
		})

	})

})
