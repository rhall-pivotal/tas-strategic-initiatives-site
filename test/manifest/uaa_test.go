package manifest_test

import (
	"strings"

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

	Describe("database connection", func() {
		Context("when PAS Database is selected", func() {
			var (
				inputProperties map[string]interface{}
			)

			BeforeEach(func() {
				inputProperties = map[string]interface{}{}
			})

			Context("and the PAS database is set to internal", func() {
				It("disables TLS to the internal database", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
					Expect(err).NotTo(HaveOccurred())

					dbAddress, err := job.Property("uaadb/address")
					Expect(err).NotTo(HaveOccurred())
					Expect(dbAddress).To(Equal("mysql.service.cf.internal"))

					tlsEnabled, err := job.Property("uaadb/tls")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(Equal("disabled"))

					caCerts, err := job.Property("uaa/ca_certs")
					Expect(err).NotTo(HaveOccurred())
					Expect(caCerts).To(HaveLen(1)) // OpsMgr root CA
				})

				Context("and TLS checkbox is checked", func() {
					BeforeEach(func() {
						inputProperties = map[string]interface{}{".properties.enable_tls_to_internal_pxc": true}
					})

					It("configures TLS to the internal database", func() {
						manifest, err := product.RenderManifest(inputProperties)
						Expect(err).NotTo(HaveOccurred())

						job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
						Expect(err).NotTo(HaveOccurred())

						tlsEnabled, err := job.Property("uaadb/tls")
						Expect(err).NotTo(HaveOccurred())
						Expect(tlsEnabled).To(Equal("enabled"))

						tlsProtocols, err := job.Property("uaadb/tls_protocols")
						Expect(err).NotTo(HaveOccurred())
						Expect(tlsProtocols).To(Equal("TLSv1.2"))

						caCerts, err := job.Property("uaa/ca_certs")
						Expect(err).NotTo(HaveOccurred())
						Expect(caCerts).To(HaveLen(1)) // OpsMgr root CA
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
						".properties.system_database.external.uaa_username":                 "some-user",
						".properties.system_database.external.uaa_password":                 map[string]interface{}{"secret": "some-password"},
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
						".properties.system_database.external.nfsvolume_username":           "nfsvolume_username",
						".properties.system_database.external.nfsvolume_password":           map[string]interface{}{"secret": "nfsvolume_password"},
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

				It("configures UAA to talk to external PAS DB", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
					Expect(err).NotTo(HaveOccurred())

					dbAddress, err := job.Property("uaadb/address")
					Expect(err).NotTo(HaveOccurred())
					Expect(dbAddress).To(Equal("foo.bar"))

					dbPort, err := job.Property("uaadb/port")
					Expect(err).NotTo(HaveOccurred())
					Expect(dbPort).To(Equal(5432))

					tlsEnabled, err := job.Property("uaadb/tls")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(Equal("disabled"))

					username, err := job.Property("uaadb/roles/0/name")
					Expect(err).NotTo(HaveOccurred())
					Expect(username).To(Equal("some-user"))

					password, err := job.Property("uaadb/roles/0/password")
					Expect(err).NotTo(HaveOccurred())
					Expect(password).NotTo(BeEmpty())

					certs, err := job.Property("uaa/ca_certs")
					Expect(err).NotTo(HaveOccurred())
					Expect(certs).To(HaveLen(2)) // OpsMgr root CA
					// UAA team told us that it's ok if this second entry is an empty string,
					// but they would fail if it was the string literal "nil"
					Expect(certs).To(ContainElement(""))
				})

				It("configures UAA to talk to DB using TLS if PAS CA cert is provided", func() {
					inputProperties[".properties.system_database.external.ca_cert"] = "some-cert"
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
					Expect(err).NotTo(HaveOccurred())

					tlsEnabled, err := job.Property("uaadb/tls")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(Equal("enabled_skip_hostname_validation"))

					caCerts, err := job.Property("uaa/ca_certs")
					Expect(err).NotTo(HaveOccurred())
					Expect(caCerts).To(HaveLen(2)) // other is OpsMgr root CA
					Expect(caCerts).To(ContainElement("some-cert"))
				})
			})
		})

		Context("when External is selected", func() {
			var inputProperties map[string]interface{}
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.uaa_database":                       "external",
					".properties.uaa_database.external.host":         "the-host",
					".properties.uaa_database.external.port":         999,
					".properties.uaa_database.external.uaa_username": "the-user",
					".properties.uaa_database.external.uaa_password": map[string]interface{}{"secret": "the-uaa-db-password"},
				}
			})

			It("configures the database", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
				Expect(err).NotTo(HaveOccurred())

				prop, err := job.Property("uaadb/db_scheme")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal("mysql"))

				prop, err = job.Property("uaadb/address")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal("the-host"))

				prop, err = job.Property("uaadb/port")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal(999))

				prop, err = job.Property("uaadb/roles/tag=admin/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal("the-user"))

				prop, err = job.Property("uaadb/roles/tag=admin/password")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(ContainSubstring("uaa_database/external/uaa_password.value"))

				prop, err = job.Property("uaadb/tls")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal("disabled"))

				prop, err = job.Property("uaadb/tls_protocols")
				Expect(err).NotTo(HaveOccurred())
				Expect(prop).To(Equal("TLSv1.2"))

				certs, err := job.Property("uaa/ca_certs")
				Expect(err).NotTo(HaveOccurred())
				Expect(certs).To(HaveLen(2)) // OpsMgr root CA
				// UAA team told us that it's ok if this second entry is an empty string,
				// but they would fail if it was the string literal "nil"
				Expect(certs).To(ContainElement(""))
			})

			Context("when a ca cert is provided", func() {
				BeforeEach(func() {
					inputProperties[".properties.uaa_database.external.ca_cert"] = "the-cert"
				})
				It("configures the database", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
					Expect(err).NotTo(HaveOccurred())

					certs, err := job.Property("uaa/ca_certs")
					Expect(err).NotTo(HaveOccurred())
					Expect(certs).To(HaveLen(2))
					Expect(certs).To(ContainElement("the-cert"))

					tlsEnabled, err := job.Property("uaadb/tls")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(Equal("enabled_skip_hostname_validation"))

					tlsProtocols, err := job.Property("uaadb/tls_protocols")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsProtocols).To(Equal("TLSv1.2"))
				})
			})
		})
	})

	Describe("route registration", func() {
		It("tags the emitted metrics", func() {
			manifest, err := product.RenderManifest(nil)
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

	Describe("Encoding compatibility mode", func() {
		It("defaults to on", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			enableUriEncodingCompatibilityMode, err := uaa.Property("uaa/authentication/enable_uri_encoding_compatibility_mode")
			Expect(err).NotTo(HaveOccurred())
			Expect(enableUriEncodingCompatibilityMode).To(BeTrue())
		})

		It("can be turned off", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".uaa.enable_uri_encoding_compatibility_mode": false,
			})

			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			enableUriEncodingCompatibilityMode, err := uaa.Property("uaa/authentication/enable_uri_encoding_compatibility_mode")
			Expect(err).NotTo(HaveOccurred())
			Expect(enableUriEncodingCompatibilityMode).To(BeFalse())
		})
	})

	Context("LDAP", func() {
		Describe("LDAP Group MaxSearchDepth", func() {
			var inputProperties map[string]interface{}

			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.uaa":          "ldap",
					".properties.uaa.ldap.url": "ldap://example.com",
					".properties.uaa.ldap.credentials": map[string]string{
						"identity": "blah",
						"password": "blah1234",
					},
					".properties.uaa.ldap.search_base":         "ou=Groups",
					".properties.uaa.ldap.search_filter":       "cn={0}",
					".properties.uaa.ldap.mail_attribute_name": "mail",
					".properties.uaa.ldap.ldap_referrals":      "follow",
					".properties.uaa.ldap.group_search_filter": "member={0}",
				}
			})

			It("defaults to 1 if not configured", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
				Expect(err).NotTo(HaveOccurred())

				ldapMaxSearchDepth, err := uaa.Property("uaa/ldap/groups/maxSearchDepth")
				Expect(err).ToNot(HaveOccurred())
				Expect(ldapMaxSearchDepth).To(Equal(1))
			})

			It("can be configured", func() {
				inputProperties[".properties.uaa.ldap.group_max_search_depth"] = "5"
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
				Expect(err).NotTo(HaveOccurred())

				ldapMaxSearchDepth, err := uaa.Property("uaa/ldap/groups/maxSearchDepth")
				Expect(err).ToNot(HaveOccurred())
				Expect(ldapMaxSearchDepth).To(Equal(5))
			})

			It("rejects values greater than the allowed maximum", func() {
				inputProperties[".properties.uaa.ldap.group_max_search_depth"] = "11"

				_, err := product.RenderManifest(inputProperties)
				Expect(err).To(HaveOccurred())
			})

			It("rejects values less than the allowed minimum", func() {
				inputProperties[".properties.uaa.ldap.group_max_search_depth"] = "0"

				_, err := product.RenderManifest(inputProperties)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("BPM", func() {
		It("co-locates and enables the BPM job with all diego jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "bpm")
			Expect(err).NotTo(HaveOccurred())

			manifestJob, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			bpmEnabled, err := manifestJob.Property("bpm/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(bpmEnabled).To(BeTrue())
		})
	})

	Context("Clients", func() {
		It("configures uaa clients", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			By("providing apps_metrics the expected permission scopes", func() {
				appMetricsScopes, err := uaa.Property("uaa/clients/apps_metrics/scope")
				Expect(err).ToNot(HaveOccurred())
				Expect(appMetricsScopes).To(Equal("cloud_controller.admin,cloud_controller.read,metrics.read,cloud_controller.admin_read_only"))
			})

			By("providing apps_metrics has the expected redirect uri", func() {
				appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics/redirect-uri")
				Expect(err).ToNot(HaveOccurred())
				Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com/**,https://metrics-previous.sys.example.com/**"))
			})

			By("providing apps_metrics_processing  the expected permission scopes", func() {
				appMetricsProcessingScopes, err := uaa.Property("uaa/clients/apps_metrics_processing/scope")
				Expect(err).ToNot(HaveOccurred())
				Expect(appMetricsProcessingScopes).To(Equal("openid,oauth.approvals,doppler.firehose,cloud_controller.admin,cloud_controller.admin_read_only"))
			})

			By("providing apps_metrics_processing the expected redirect uri", func() {
				appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics_processing/redirect-uri")
				Expect(err).ToNot(HaveOccurred())
				Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com,https://metrics-previous.sys.example.com"))
			})

			By("providing catalyst the expected permission scopes", func() {
				catalystScopes, err := uaa.Property("uaa/clients/catalyst/scope")
				Expect(err).ToNot(HaveOccurred())
				Expect(catalystScopes).To(Equal("uaa.user,cloud_controller.read,cloud_controller.admin_read_only,logs.admin,doppler.firehose"))
			})

			By("providing catalyst has the expected redirect uri", func() {
				catalystRedirectUri, err := uaa.Property("uaa/clients/catalyst/redirect-uri")
				Expect(err).ToNot(HaveOccurred())
				Expect(catalystRedirectUri).To(Equal("" +
					"https://metrics.sys.example.com/**," +
					"https://appwatch.sys.example.com/**," +
					"https://appmetrics.sys.example.com/**," +
					"https://applogs.sys.example.com/**"))
			})

			By("providing apps_manager_js client the expected scopes", func() {
				rawScopes, err := uaa.Property("uaa/clients/apps_manager_js/scope")
				Expect(err).ToNot(HaveOccurred())

				scopes := strings.Split(rawScopes.(string), ",")
				Expect(scopes).To(ContainElement("network.write"))
				Expect(scopes).To(ContainElement("network.admin"))

				autoapproveList, err := uaa.Property("uaa/clients/apps_manager_js/autoapprove")
				Expect(err).ToNot(HaveOccurred())

				Expect(autoapproveList).To(ContainElement("network.write"))
				Expect(autoapproveList).To(ContainElement("network.admin"))
			})

			By("providing credhub_admin_client the expected scopes", func() {
				id, err := uaa.Property("uaa/clients/credhub_admin_client/id")
				Expect(err).ToNot(HaveOccurred())
				Expect(id).To(Equal("credhub_admin_client"))

				rawAuthorities, err := uaa.Property("uaa/clients/credhub_admin_client/authorities")
				Expect(err).ToNot(HaveOccurred())

				authorities := strings.Split(rawAuthorities.(string), ",")
				Expect(authorities).To(ConsistOf([]string{"credhub.read", "credhub.write"}))

				authorizedGrantTypes, err := uaa.Property("uaa/clients/credhub_admin_client/authorized-grant-types")
				Expect(err).ToNot(HaveOccurred())
				Expect(authorizedGrantTypes).To(Equal("client_credentials"))
			})

			By("providing tile_installer with the right properties", func() {
				id, err := uaa.Property("uaa/clients/tile_installer/id")
				Expect(err).ToNot(HaveOccurred())
				Expect(id).To(Equal("tile_installer"))

				rawAuthorities, err := uaa.Property("uaa/clients/tile_installer/authorities")
				Expect(err).ToNot(HaveOccurred())

				authorities := strings.Split(rawAuthorities.(string), ",")
				Expect(authorities).To(ConsistOf([]string{"cloud_controller.admin", "clients.admin", "credhub.read", "credhub.write"}))

				authorizedGrantTypes, err := uaa.Property("uaa/clients/tile_installer/authorized-grant-types")
				Expect(err).ToNot(HaveOccurred())
				Expect(authorizedGrantTypes).To(Equal("client_credentials"))

				accessTokenValidity, err := uaa.Property("uaa/clients/tile_installer/access-token-validity")
				Expect(err).ToNot(HaveOccurred())
				Expect(accessTokenValidity).To(Equal(3600))

				override, err := uaa.Property("uaa/clients/tile_installer/override")
				Expect(err).ToNot(HaveOccurred())
				Expect(override).To(BeTrue())
			})

			By("allowing users to login to usage service with token", func() {
				rawScopes, err := uaa.Property("uaa/clients/cf/scope")
				Expect(err).ToNot(HaveOccurred())
				scopes := strings.Split(rawScopes.(string), ",")
				Expect(scopes).To(ContainElement("usage_service.audit"))

				rawGroups, err := uaa.Property("uaa/scim/groups")
				groups := rawGroups.(map[interface{}]interface{})
				Expect(err).ToNot(HaveOccurred())
				Expect(groups).To(HaveKeyWithValue("usage_service.audit", "View reports for the Usage Service"))
			})

			By("metric registrar smoke test using uaa client", func() {
				rawAuthorities, err := uaa.Property("uaa/clients/metric_registrar_smoke_test/authorities")
				Expect(err).ToNot(HaveOccurred())
				authorities := strings.Split(rawAuthorities.(string), ",")
				Expect(authorities).To(ContainElement("cloud_controller.admin"))

			})
		})
	})
})
