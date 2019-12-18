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

	Describe("internal-provider", func() {
		Context("when there is a single internal key", func() {
			It("configures credhub with the key", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				keys, err := credhub.Property("credhub/encryption/keys")
				Expect(err).ToNot(HaveOccurred())
				Expect(keys).To(HaveLen(2))

				key := keys.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(key["provider_name"]).To(Equal("internal-provider"))
				Expect(key["key_properties"]).To(HaveKeyWithValue("encryption_password",
					ContainSubstring("credhub_internal_provider_keys/0/key.value")))
				Expect(key["active"]).To(BeTrue())

				providers, err := credhub.Property("credhub/encryption/providers")

				Expect(providers).To(HaveLen(2))
				Expect(err).ToNot(HaveOccurred())

				internalProvider := providers.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(internalProvider["name"]).To(Equal("internal-provider"))
				Expect(internalProvider["type"]).To(Equal("internal"))
			})
		})

		Context("when there are multiple internal keys", func() {
			It("configures credhub with multiple keys", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.credhub_internal_provider_keys": []map[string]interface{}{
						{
							"key": map[string]interface{}{
								"secret": "12345678901234567890",
							},
							"name":    "KeyA",
							"primary": true,
						},
						{
							"key": map[string]interface{}{
								"secret": "09876543210987654321",
							},
							"name":    "KeyB",
							"primary": false,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				keys, err := credhub.Property("credhub/encryption/keys")
				Expect(err).ToNot(HaveOccurred())
				Expect(keys).To(HaveLen(2))

				firstKey := keys.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(firstKey["provider_name"]).To(Equal("internal-provider"))
				Expect(firstKey["key_properties"]).To(HaveKeyWithValue("encryption_password",
					ContainSubstring("credhub_internal_provider_keys/0/key.value")))
				Expect(firstKey["active"]).To(BeTrue())

				secondKey := keys.([]interface{})[0].([]interface{})[1].(map[interface{}]interface{})
				Expect(secondKey["provider_name"]).To(Equal("internal-provider"))
				Expect(secondKey["key_properties"]).To(HaveKeyWithValue("encryption_password",
					ContainSubstring("credhub_internal_provider_keys/1/key.value")))
				Expect(secondKey["active"]).To(BeFalse())

				providers, err := credhub.Property("credhub/encryption/providers")
				Expect(providers).To(HaveLen(2))
				Expect(err).ToNot(HaveOccurred())

				internalProvider := providers.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(internalProvider["name"]).To(Equal("internal-provider"))
				Expect(internalProvider["type"]).To(Equal("internal"))
			})
		})
	})

	Describe("kms provider", func() {
		It("configures credhub with multiple keys", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.credhub_kms_providers": []map[string]interface{}{
					{
						"instance_name": "some-kms-instance-name-provider",
						"endpoint":      "some-endpoint",
						"primary":       true,
					},
				},
				".properties.credhub_internal_provider_keys": []map[string]interface{}{
					{
						"key": map[string]interface{}{
							"secret": "12345678901234567890",
						},
						"name":    "KeyA",
						"primary": false,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			keys, err := credhub.Property("credhub/encryption/keys")
			Expect(err).ToNot(HaveOccurred())
			Expect(keys).To(HaveLen(2))

			internalProviderKey := keys.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
			Expect(internalProviderKey["provider_name"]).To(Equal("internal-provider"))
			Expect(internalProviderKey["key_properties"]).To(HaveKeyWithValue("encryption_password",
				ContainSubstring("credhub_internal_provider_keys/0/key.value")))
			Expect(internalProviderKey["active"]).To(BeFalse())

			kmsPluginProviderKey := keys.([]interface{})[1].([]interface{})[0].(map[interface{}]interface{})
			Expect(kmsPluginProviderKey["provider_name"]).To(Equal("some-kms-instance-name-provider"))
			Expect(kmsPluginProviderKey["key_properties"]).To(HaveKeyWithValue("encryption_key_name",
				ContainSubstring("kms-plugin-key-name")))
			Expect(kmsPluginProviderKey["active"]).To(BeTrue())

			providers, err := credhub.Property("credhub/encryption/providers")
			Expect(providers).To(HaveLen(2))
			Expect(err).ToNot(HaveOccurred())

			internalProvider := providers.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
			Expect(internalProvider["name"]).To(Equal("internal-provider"))
			Expect(internalProvider["type"]).To(Equal("internal"))

			kmsProvider := providers.([]interface{})[1].([]interface{})[0].(map[interface{}]interface{})
			Expect(kmsProvider["name"]).To(Equal("some-kms-instance-name-provider"))
			Expect(kmsProvider["type"]).To(Equal("kms-plugin"))
			Expect(kmsProvider["connection_properties"]).To(HaveKeyWithValue("endpoint", "some-endpoint"))
			Expect(kmsProvider["connection_properties"]).To(HaveKeyWithValue("host", "credhub-kms"))
			Expect(kmsProvider["connection_properties"]).To(HaveKeyWithValue("ca", "((/services/tls_ca))"))
		})
	})

	Describe("encryption keys", func() {
		Context("when there is an additional HSM key set as primary", func() {
			It("configures credhub with the keys, with the HSM key marked as active", func() {
				fakeClientKeypair := generateTLSKeypair("some-hsm-client")
				fakeServerKeypair := generateTLSKeypair("some-hsm-host")
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.credhub_internal_provider_keys": []map[string]interface{}{
						{
							"key": map[string]interface{}{
								"secret": "12345678901234567890",
							},
							"name":    "KeyA",
							"primary": false,
						},
						{
							"key": map[string]interface{}{
								"secret": "98765432109876543210",
							},
							"name":    "KeyB",
							"primary": false,
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
							"certificate":             fakeServerKeypair.Certificate,
							"partition_serial_number": "some-hsm-partition-serial",
							"port": 9999,
						},
					},
					".properties.credhub_hsm_provider_encryption_keys": []map[string]interface{}{
						{
							"name":    "some-hsm-key-name",
							"primary": true,
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				keys, err := credhub.Property("credhub/encryption/keys")
				Expect(err).ToNot(HaveOccurred())
				Expect(keys).To(HaveLen(3))

				firstInternalKey := keys.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(firstInternalKey["provider_name"]).To(Equal("internal-provider"))
				Expect(firstInternalKey["key_properties"]).To(HaveKeyWithValue("encryption_password", ContainSubstring("credhub_internal_provider_keys/0/key.value")))
				Expect(firstInternalKey["active"]).To(BeFalse())

				secondInternalKey := keys.([]interface{})[0].([]interface{})[1].(map[interface{}]interface{})
				Expect(secondInternalKey["provider_name"]).To(Equal("internal-provider"))
				Expect(secondInternalKey["key_properties"]).To(HaveKeyWithValue("encryption_password", ContainSubstring("credhub_internal_provider_keys/1/key.value")))
				Expect(secondInternalKey["active"]).To(BeFalse())

				hsmKey := keys.([]interface{})[2].([]interface{})[0].(map[interface{}]interface{})
				Expect(hsmKey["provider_name"]).To(Equal("hsm-provider"))
				Expect(hsmKey["key_properties"]).To(HaveKeyWithValue("encryption_key_name", "some-hsm-key-name"))
				Expect(hsmKey["active"]).To(BeTrue())

				providers, err := credhub.Property("credhub/encryption/providers")
				Expect(err).ToNot(HaveOccurred())
				Expect(providers).To(HaveLen(3))

				internalProvider := providers.([]interface{})[0].([]interface{})[0].(map[interface{}]interface{})
				Expect(internalProvider["name"]).To(Equal("internal-provider"))
				Expect(internalProvider["type"]).To(Equal("internal"))

				hsmProvider := providers.([]interface{})[2].(map[interface{}]interface{})
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

	Describe("database configuration", func() {
		Context("when PAS Database is selected", func() {
			Context("and the PAS database is set to internal", func() {
				It("configures credhub and bbr-credhubdb to talk to mysql without tls", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
					Expect(err).NotTo(HaveOccurred())

					requireTLS, err := credhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeFalse())

					ca, err := credhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).To(BeNil())

					bbrCredhub, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-credhubdb")
					Expect(err).NotTo(HaveOccurred())

					requireTLS, err = bbrCredhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeFalse())

					ca, err = bbrCredhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).To(BeNil())
				})

				Context("when tls checkbox is checked", func() {
					It("configures credhubdb to use tls", func() {
						manifest, err := product.RenderManifest(map[string]interface{}{".properties.enable_tls_to_internal_pxc": true})
						Expect(err).NotTo(HaveOccurred())

						credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
						Expect(err).NotTo(HaveOccurred())

						requireTLS, err := credhub.Property("credhub/data_storage/require_tls")
						Expect(err).ToNot(HaveOccurred())
						Expect(requireTLS).To(BeTrue())

						ca, err := credhub.Property("credhub/data_storage/tls_ca")
						Expect(err).ToNot(HaveOccurred())
						Expect(ca).NotTo(BeEmpty())

						bbrCredhub, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-credhubdb")
						Expect(err).NotTo(HaveOccurred())

						requireTLS, err = bbrCredhub.Property("credhub/data_storage/require_tls")
						Expect(err).ToNot(HaveOccurred())
						Expect(requireTLS).To(BeTrue())

						ca, err = bbrCredhub.Property("credhub/data_storage/tls_ca")
						Expect(err).ToNot(HaveOccurred())
						Expect(ca).NotTo(BeEmpty())
					})
				})
			})

			Context("and the PAS database is set to external", func() {
				var inputProperties map[string]interface{}

				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".properties.system_database":                                       "external",
						".properties.system_database.external.host":                         "foo.bar",
						".properties.system_database.external.port":                         5432,
						".properties.system_database.external.credhub_username":             "some-user",
						".properties.system_database.external.credhub_password":             map[string]interface{}{"secret": "some-password"},
						".properties.system_database.external.app_usage_service_username":   "app_usage_service_username",
						".properties.system_database.external.app_usage_service_password":   map[string]interface{}{"secret": "app_usage_service_password"},
						".properties.system_database.external.autoscale_username":           "autoscale_username",
						".properties.system_database.external.autoscale_password":           map[string]interface{}{"secret": "autoscale_password"},
						".properties.system_database.external.ccdb_username":                "ccdb_username",
						".properties.system_database.external.ccdb_password":                map[string]interface{}{"secret": "ccdb_password"},
						".properties.system_database.external.diego_username":               "diego_username",
						".properties.system_database.external.diego_password":               map[string]interface{}{"secret": "diego_password"},
						".properties.system_database.external.locket_username":              "locket_username",
						".properties.system_database.external.locket_password":              map[string]interface{}{"secret": "locket_password"},
						".properties.system_database.external.networkpolicyserver_username": "networkpolicyserver_username",
						".properties.system_database.external.networkpolicyserver_password": map[string]interface{}{"secret": "networkpolicyserver_password"},
						".properties.system_database.external.notifications_username":       "notifications_username",
						".properties.system_database.external.notifications_password":       map[string]interface{}{"secret": "notifications_password"},
						".properties.system_database.external.account_username":             "account_username",
						".properties.system_database.external.account_password":             map[string]interface{}{"secret": "account_password"},
						".properties.system_database.external.routing_username":             "routing_username",
						".properties.system_database.external.routing_password":             map[string]interface{}{"secret": "routing_password"},
						".properties.system_database.external.silk_username":                "silk_username",
						".properties.system_database.external.silk_password":                map[string]interface{}{"secret": "silk_password"},
					}
				})

				It("configures credhub and bbr-credhubdb to talk to external PAS DB", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
					Expect(err).NotTo(HaveOccurred())

					host, err := credhub.Property("credhub/data_storage/host")
					Expect(err).ToNot(HaveOccurred())
					Expect(host).To(Equal("foo.bar"))

					port, err := credhub.Property("credhub/data_storage/port")
					Expect(err).ToNot(HaveOccurred())
					Expect(port).To(Equal(5432))

					dbName, err := credhub.Property("credhub/data_storage/database")
					Expect(err).ToNot(HaveOccurred())
					Expect(dbName).To(Equal("credhub"))

					dbUsername, err := credhub.Property("credhub/data_storage/username")
					Expect(err).ToNot(HaveOccurred())
					Expect(dbUsername).To(Equal("some-user"))

					dbPassword, err := credhub.Property("credhub/data_storage/password")
					Expect(err).ToNot(HaveOccurred())
					Expect(dbPassword).NotTo(BeNil())

					requireTLS, err := credhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeFalse())

					ca, err := credhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).To(BeNil())

					bbrCredhub, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-credhubdb")
					Expect(err).NotTo(HaveOccurred())

					requireTLS, err = bbrCredhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeFalse())

					ca, err = bbrCredhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).To(BeNil())
				})

				It("configures credhub and bbr-credhubdb to use TLS if PAS CA cert is provided", func() {
					inputProperties[".properties.system_database.external.ca_cert"] = "some-cert"
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
					Expect(err).NotTo(HaveOccurred())

					requireTLS, err := credhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeTrue())

					ca, err := credhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).NotTo(BeEmpty())

					bbrCredhub, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-credhubdb")
					Expect(err).NotTo(HaveOccurred())

					requireTLS, err = bbrCredhub.Property("credhub/data_storage/require_tls")
					Expect(err).ToNot(HaveOccurred())
					Expect(requireTLS).To(BeTrue())

					ca, err = bbrCredhub.Property("credhub/data_storage/tls_ca")
					Expect(err).ToNot(HaveOccurred())
					Expect(ca).NotTo(BeEmpty())
				})
			})
		})

		Context("when External is selected", func() {
			inputProperties := map[string]interface{}{
				".properties.credhub_database":                   "external",
				".properties.credhub_database.external.tls_ca":   "fake-ca",
				".properties.credhub_database.external.host":     "cred.foo.bar",
				".properties.credhub_database.external.port":     "2345",
				".properties.credhub_database.external.username": "credhub_username",
				".properties.credhub_database.external.password": map[string]interface{}{"secret": "credhub_password"},
			}

			It("it requires TLS", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
				Expect(err).NotTo(HaveOccurred())

				requireTLS, err := credhub.Property("credhub/data_storage/require_tls")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireTLS).To(BeTrue())

				ca, err := credhub.Property("credhub/data_storage/tls_ca")
				Expect(err).NotTo(HaveOccurred())
				Expect(ca).NotTo(BeEmpty())
			})
			Context("disable_hostname_verification", func() {
				It("disables hostname_verification when selected", func() {
					inputProperties[".properties.credhub_database.external.disable_hostname_verification"] = true
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
					Expect(err).NotTo(HaveOccurred())

					hostnameVerification, err := credhub.Property("credhub/data_storage/hostname_verification/enabled")
					Expect(err).NotTo(HaveOccurred())
					Expect(hostnameVerification).To(BeFalse())
				})
				It("does not disable hostname_verification by default", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
					Expect(err).NotTo(HaveOccurred())

					hostnameVerification, err := credhub.Property("credhub/data_storage/hostname_verification/enabled")
					Expect(err).NotTo(HaveOccurred())
					Expect(hostnameVerification).To(BeTrue())
				})
			})
		})
	})

	Describe("permissions", func() {
		It("provides uaa operations rights", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			permissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).ToNot(HaveOccurred())

			By("granting permissions to the credhub-service-broker tile")
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

			By("granting permissions to the services_credhub_client")
			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/c/*",
				"actors":     []interface{}{"uaa-client:services_credhub_client"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))

			By("granting permissions to the cloud controller to read service key credentials")
			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/*",
				"actors":     []interface{}{"uaa-client:cc_service_key_client"},
				"operations": []interface{}{"read"},
			}))

			By("granting permissions to the credhub_admin_client")
			Expect(permissions).To(ContainElement(map[interface{}]interface{}{
				"path":       "/*",
				"actors":     []interface{}{"uaa-client:credhub_admin_client"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
			}))
		})
	})
})
