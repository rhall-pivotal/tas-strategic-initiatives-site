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
				fakeClientKeypair := generateTLSKeypair("some-hsm-client")
				fakeServerKeypair := generateTLSKeypair("some-hsm-host")
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
					".properties.credhub_hsm_provider_client_certificate": map[string]interface{}{
						"cert_pem":        fakeClientKeypair.Certificate,
						"private_key_pem": fakeClientKeypair.PrivateKey,
					},
					".properties.credhub_hsm_provider_partition": "some-hsm-partition",
					".properties.credhub_hsm_provider_partition_password": map[string]interface{}{
						"secret": "some-hsm-partition-password",
					},
					".properties.credhub_hsm_provider_servers": []map[string]interface{}{
						{
							"host_address":            "some-hsm-host",
							"hsm_certificate":         fakeServerKeypair.Certificate,
							"partition_serial_number": "some-hsm-partition-serial",
							"port": 9999,
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

				providers, err := credhub.Property("credhub/encryption/providers")
				Expect(err).ToNot(HaveOccurred())

				internalProvider := providers.([]interface{})[0].(map[interface{}]interface{})
				Expect(internalProvider["name"]).To(Equal("internal-provider"))
				Expect(internalProvider["type"]).To(Equal("internal"))

				hsmProvider := providers.([]interface{})[1].(map[interface{}]interface{})
				Expect(hsmProvider["name"]).To(Equal("hsm-provider"))
				Expect(hsmProvider["type"]).To(Equal("hsm"))

				hsmConnectionProperties := hsmProvider["connection_properties"].(map[interface{}]interface{})
				Expect(hsmConnectionProperties["partition"]).To(Equal("some-hsm-partition"))
				Expect(hsmConnectionProperties["partition_password"]).NotTo(BeEmpty())
				Expect(hsmConnectionProperties["client_certificate"]).NotTo(BeEmpty())
				Expect(hsmConnectionProperties["client_key"]).NotTo(BeEmpty())

				hsmServer := hsmConnectionProperties["servers"].([]interface{})[0].(map[interface{}]interface{})
				Expect(hsmServer["certificate"]).NotTo(BeEmpty())
				Expect(hsmServer["host"]).To(Equal("some-hsm-host"))
				Expect(hsmServer["partition_serial_number"]).To(Equal("some-hsm-partition-serial"))
				Expect(hsmServer["port"]).To(Equal(9999))
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

		It("grants permission to the cloud controller to read service key credentials", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			permissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).ToNot(HaveOccurred())

			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/*",
				"actors":     []interface{}{"uaa-client:cc_service_key_client"},
				"operations": []interface{}{"read"},
			}))
		})

		It("grants permissions to the credhub_admin_client", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			permissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).ToNot(HaveOccurred())

			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/*",
				"actors":     []interface{}{"uaa-client:credhub_admin_client"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))
		})

	})

})
