package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Networking", func() {
	Describe("Container networking", func() {
		var (
			inputProperties         map[string]interface{}
			cellInstanceGroup       string
			controllerInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "ert" {
				cellInstanceGroup = "diego_cell"
				controllerInstanceGroup = "diego_database"
			} else {
				cellInstanceGroup = "compute"
				controllerInstanceGroup = "control"
			}
			inputProperties = map[string]interface{}{}
		})

		Describe("policy server", func() {
			It("uses the correct database host", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				host, err := job.Property("database/host")
				Expect(err).NotTo(HaveOccurred())
				Expect(host).To(Equal("mysql.service.cf.internal"))

				databaseLink, err := job.Path("/consumes/database")
				Expect(err).NotTo(HaveOccurred())
				Expect(databaseLink).To(Equal("nil"))
			})

			Context("when the operator does not set a limit for policy server open database connections", func() {
				It("configures jobs with default values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(200))
				})
			})

			Context("when the user specifies custom values for policy server max open database connections", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".properties.networkpolicyserver_database_max_open_connections": 300,
					}
				})

				It("configures jobs with user provided values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(300))
				})

				Context("when the policy server max open DB connections is out of range", func() {
					BeforeEach(func() {
						inputProperties = map[string]interface{}{
							".properties.networkpolicyserver_database_max_open_connections": 0,
						}
					})

					It("returns an error", func() {
						_, err := product.RenderManifest(inputProperties)
						Expect(err.Error()).To(ContainSubstring("Value must be greater than or equal to 1"))
					})
				})
			})

			It("disables TLS to the internal database by default", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				tlsEnabled, err := job.Property("database/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(tlsEnabled).To(BeFalse())

				caCert, err := job.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(BeNil())
			})

			Context("when the TLS checkbox is checked", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{".properties.enable_tls_to_internal_pxc": true}
				})

				It("configures TLS to the internal database", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
					Expect(err).NotTo(HaveOccurred())

					tlsEnabled, err := job.Property("database/require_ssl")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(BeTrue())

					caCert, err := job.Property("database/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(caCert).NotTo(BeEmpty())
				})
			})
		})

		Describe("policy server internal", func() {
			Context("when experimental dynamic egress enforcement is enabled", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".properties.experimental_dynamic_egress_enforcement": true,
					}
				})

				It("enables experimental dynamic egress policy", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server-internal")
					Expect(err).NotTo(HaveOccurred())

					enabled, err := job.Property("enforce_experimental_dynamic_egress_policies")
					Expect(err).NotTo(HaveOccurred())
					Expect(enabled).To(BeTrue())
				})
			})

			Context("when the operator does not set a limit for policy server internal open database connections", func() {
				It("configures jobs with default values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server-internal")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(200))
				})
			})

			Context("when the user specifies custom values for policy server internal max open database connections", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".properties.networkpolicyserverinternal_database_max_open_connections": 300,
					}
				})

				It("configures jobs with user provided values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server-internal")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(300))
				})

				Context("when the policy server internal max open DB connections is out of range", func() {
					BeforeEach(func() {
						inputProperties = map[string]interface{}{
							".properties.networkpolicyserverinternal_database_max_open_connections": 0,
						}
					})

					It("returns an error", func() {
						_, err := product.RenderManifest(inputProperties)
						Expect(err.Error()).To(ContainSubstring("Value must be greater than or equal to 1"))
					})
				})
			})
		})

		Context("when the operator configures database connection timeout for CNI plugin", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.cf_networking_database_connection_timeout": 250,
				}
			})

			It("sets the manifest database connection timeout properties for the cf networking jobs to be 250", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				policyServerJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				policyServerConnectTimeoutSeconds, err := policyServerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerConnectTimeoutSeconds).To(Equal(250))

				policyServerInternalJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server-internal")
				Expect(err).NotTo(HaveOccurred())

				policyServerInternalConnectTimeoutSeconds, err := policyServerInternalJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(policyServerInternalConnectTimeoutSeconds).To(Equal(250))

				silkControllerJob, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				silkControllerConnectTimeoutSeconds, err := silkControllerJob.Property("database/connect_timeout_seconds")
				Expect(err).NotTo(HaveOccurred())

				Expect(silkControllerConnectTimeoutSeconds).To(Equal(250))
			})
		})

		Context("when Silk is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{}
			})

			It("configures the cni_config_dir and cni_plugin_dir", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/silk-cni/config/cni"))

				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniPluginDir).To(Equal("/var/vcap/packages/silk-cni/bin"))
			})

			It("disables TLS to the internal database by default", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				tlsEnabled, err := job.Property("database/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(tlsEnabled).To(BeFalse())

				caCert, err := job.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(BeNil())
			})

			Context("when the TLS checkbox is checked", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{".properties.enable_tls_to_internal_pxc": true}
				})

				It("configures TLS to the internal database", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
					Expect(err).NotTo(HaveOccurred())

					tlsEnabled, err := job.Property("database/require_ssl")
					Expect(err).NotTo(HaveOccurred())
					Expect(tlsEnabled).To(BeTrue())

					caCert, err := job.Property("database/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(caCert).NotTo(BeEmpty())
				})
			})

			It("uses the correct database host", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				host, err := job.Property("database/host")
				Expect(err).NotTo(HaveOccurred())
				Expect(host).To(Equal("mysql.service.cf.internal"))

				databaseLink, err := job.Path("/consumes/database")
				Expect(err).NotTo(HaveOccurred())
				Expect(databaseLink).To(Equal("nil"))
			})

			Context("when the operator does not set a limit for silk-controller open database connections", func() {
				It("configures jobs with default values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(200))
				})
			})

			Context("when the user specifies custom values for silk-controller max open database connections", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".properties.silk_database_max_open_connections": 300,
					}
				})

				It("configures jobs with user provided values", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "silk-controller")
					Expect(err).NotTo(HaveOccurred())

					maxOpenConnections, err := job.Property("max_open_connections")
					Expect(err).NotTo(HaveOccurred())
					Expect(maxOpenConnections).To(Equal(300))
				})

				Context("when the silk-controller max open DB connections is out of range", func() {
					BeforeEach(func() {
						inputProperties = map[string]interface{}{
							".properties.silk_database_max_open_connections": 0,
						}
					})

					It("returns an error", func() {
						_, err := product.RenderManifest(inputProperties)
						Expect(err.Error()).To(ContainSubstring("Value must be greater than or equal to 1"))
					})
				})
			})

			Context("silk network policy", func() {
				It("continues to be enforced", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "vxlan-policy-agent")
					Expect(err).NotTo(HaveOccurred())

					disabled, err := job.Property("disable_container_network_policy")
					Expect(err).NotTo(HaveOccurred())
					Expect(disabled).To(BeFalse())
				})

				Context("setting is disabled", func() {
					BeforeEach(func() {
						inputProperties = map[string]interface{}{
							".properties.container_networking_interface_plugin.silk.enable_policy_enforcement": false,
						}
					})

					It("disables silk network policy enforcement", func() {
						manifest, err := product.RenderManifest(inputProperties)
						Expect(err).NotTo(HaveOccurred())

						job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "vxlan-policy-agent")
						Expect(err).NotTo(HaveOccurred())

						disabled, err := job.Property("disable_container_network_policy")
						Expect(err).NotTo(HaveOccurred())
						Expect(disabled).To(BeTrue())
					})
				})
			})
		})

		Context("when External is enabled", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.container_networking_interface_plugin": "external",
				}
			})

			It("configures the cni_config_dir", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(cellInstanceGroup, "garden-cni")
				Expect(err).NotTo(HaveOccurred())

				cniConfigDir, err := job.Property("cni_config_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniConfigDir).To(Equal("/var/vcap/jobs/cni/config/cni"))

				cniPluginDir, err := job.Property("cni_plugin_dir")
				Expect(err).NotTo(HaveOccurred())
				Expect(cniPluginDir).To(Equal("/var/vcap/packages/cni/bin"))
			})
		})
	})

	Describe("DNS search domain", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "diego_cell"
			} else {
				instanceGroup = "compute"
			}
		})

		It("configures search_domains on the garden-cni job", func() {
			inputProperties = map[string]interface{}{
				".properties.cf_networking_search_domains": "some-search-domain,another-search-domain",
			}

			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(Equal([]interface{}{
				"some-search-domain",
				"another-search-domain",
			}))
		})

		It("configures search_domains on the garden-cni job", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "garden-cni")
			Expect(err).NotTo(HaveOccurred())

			searchDomains, err := job.Property("search_domains")
			Expect(err).NotTo(HaveOccurred())

			Expect(searchDomains).To(HaveLen(0))
		})
	})

	Describe("Service Discovery For Apps", func() {
		Describe("controller", func() {
			var instanceGroup string
			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "diego_brain"
				} else {
					instanceGroup = "control"
				}
			})

			It("is deployed", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "service-discovery-controller")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Describe("cell", func() {
			var instanceGroup string
			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "diego_cell"
				} else {
					instanceGroup = "compute"
				}
			})

			It("co-locates the bosh-dns-adapter and bpm", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
				Expect(err).NotTo(HaveOccurred())

				_, err = manifest.FindInstanceGroupJob(instanceGroup, "bpm")
				Expect(err).NotTo(HaveOccurred())
			})

			It("emits internal routes", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "route_emitter")
				Expect(err).NotTo(HaveOccurred())

				enabled, err := job.Property("internal_routes/enabled")
				Expect(err).NotTo(HaveOccurred())

				Expect(enabled).To(BeTrue())
			})

			Context("when internal domain is empty", func() {
				It("defaults internal domain to apps.internal", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("internal_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(ConsistOf("apps.internal"))
				})
			})

			Context("when internal domains are configured", func() {
				var (
					inputProperties map[string]interface{}
				)

				It("sets internal domains to the provided internal domains", func() {
					inputProperties = map[string]interface{}{
						".properties.cf_networking_internal_domains": []map[string]interface{}{
							{
								"name": "some-internal-domain",
							},
							{
								"name": "some-other-internal-domain",
							},
						},
					}
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "bosh-dns-adapter")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("internal_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"some-internal-domain",
						"some-other-internal-domain",
					}))
				})
			})
		})

		Describe("api", func() {
			var instanceGroup string
			BeforeEach(func() {
				if productName == "ert" {
					instanceGroup = "cloud_controller"
				} else {
					instanceGroup = "control"
				}
			})

			Context("when internal domain is empty", func() {
				It("adds apps.internal to app domains", func() {
					manifest, err := product.RenderManifest(nil)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("app_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"apps.example.com",
						[]interface{}{},
						[]interface{}{
							map[interface{}]interface{}{
								"internal": true,
								"name":     "apps.internal",
							},
						},
					}))
				})
			})

			Context("when internal domains are configured", func() {
				var (
					inputProperties map[string]interface{}
				)

				It("adds internal domains to app domains", func() {
					inputProperties = map[string]interface{}{
						".properties.cf_networking_internal_domains": []map[string]interface{}{
							{
								"name": "some-internal-domain",
							},
							{
								"name": "some-other-internal-domain",
							},
						},
					}
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("app_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"apps.example.com",
						[]interface{}{},
						[]interface{}{
							map[interface{}]interface{}{
								"name":     "some-internal-domain",
								"internal": true,
							},
							map[interface{}]interface{}{
								"name":     "some-other-internal-domain",
								"internal": true,
							},
						},
					}))
				})
			})
		})

		Describe("SSH proxy", func() {
			var instanceGroup string
			BeforeEach(func() {
				if productName == "srt" {
					instanceGroup = "control"
				} else {
					instanceGroup = "diego_brain"
				}
			})

			Context("when static IPs are set", func() {
				var inputProperties map[string]interface{}

				It("adds the static_ips", func() {
					p := ""
					if instanceGroup == "diego_brain" {
						p = ".diego_brain.static_ips"
					} else {
						p = ".control.static_ips"
					}

					inputProperties = map[string]interface{}{
						p: "0.0.0.0",
					}

					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					ips, err := manifest.Path(fmt.Sprintf("/instance_groups/name=%s/networks", instanceGroup))
					Expect(err).NotTo(HaveOccurred())

					key := ips.([]interface{})[0].(map[interface{}]interface{})
					Expect(key["static_ips"]).To(Equal([]interface{}{"0.0.0.0"}))
				})
			})
		})

		Describe("Istio", func() {
			var capiInstanceGroup string
			BeforeEach(func() {
				if productName == "srt" {
					capiInstanceGroup = "control"
				} else {
					capiInstanceGroup = "cloud_controller"
				}
			})

			Context("when it is enabled", func() {
				It("adds does not zero out istio-control, istio-router, or cc_route_syncer", func() {
					inputProperties := map[string]interface{}{
						".properties.istio": "enable",
					}

					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					instanceCount, err := manifest.Path("/instance_groups/name=istio_control/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(1))

					instanceCount, err = manifest.Path("/instance_groups/name=istio_router/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(2))

					instanceCount, err = manifest.Path("/instance_groups/name=route_syncer/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(1))
				})

				Context("when static IPs are set", func() {
					It("adds the static_ips", func() {
						inputProperties := map[string]interface{}{
							".properties.istio":        "enable",
							".istio_router.static_ips": "0.0.0.0",
						}

						manifest, err := product.RenderManifest(inputProperties)
						Expect(err).NotTo(HaveOccurred())

						ips, err := manifest.Path("/instance_groups/name=istio_router/networks")
						Expect(err).NotTo(HaveOccurred())

						key := ips.([]interface{})[0].(map[interface{}]interface{})
						Expect(key["static_ips"]).To(Equal([]interface{}{"0.0.0.0"}))
					})
				})

				Describe("when frontend TLS keypairs are configured", func() {
					It("populates the frontend TLS keypairs", func() {
						fakeCert1 := generateTLSKeypair("some-hostname")
						fakeCert2 := generateTLSKeypair("another-hostname")
						inputProperties := map[string]interface{}{
							".properties.istio": "enable",
							".properties.istio_frontend_tls_keypairs": []map[string]interface{}{
								{"name": "cert-1", "certificate": map[string]interface{}{"cert_pem": fakeCert1.Certificate, "private_key_pem": fakeCert1.PrivateKey}},
								{"name": "cert-2", "certificate": map[string]interface{}{"cert_pem": fakeCert2.Certificate, "private_key_pem": fakeCert2.PrivateKey}},
							},
						}
						manifest, err := product.RenderManifest(inputProperties)
						Expect(err).NotTo(HaveOccurred())

						copilot, err := manifest.FindInstanceGroupJob("istio_control", "copilot")
						Expect(err).NotTo(HaveOccurred())
						keyPairsInterface, err := copilot.Property("frontend_tls_keypairs")
						keyPairs := keyPairsInterface.([]interface{})
						Expect(err).NotTo(HaveOccurred())
						Expect(len(keyPairs)).To(Equal(2))
						kp := keyPairs[0].(map[interface{}]interface{})
						Expect(kp["cert_chain"]).NotTo(BeEmpty())
						Expect(kp["private_key"]).NotTo(BeEmpty())
						kp = keyPairs[1].(map[interface{}]interface{})
						Expect(kp["cert_chain"]).NotTo(BeEmpty())
						Expect(kp["private_key"]).NotTo(BeEmpty())
					})
				})

				Describe("istio_domain", func() {
					Describe("default", func() {
						var instanceGroup string
						BeforeEach(func() {
							if productName == "ert" {
								instanceGroup = "cloud_controller"
							} else {
								instanceGroup = "control"
							}
						})

						It("adds default domains to both app domains and temporary istio domains", func() {
							inputProperties := map[string]interface{}{
								".properties.istio": "enable",
							}
							manifest, err := product.RenderManifest(inputProperties)
							Expect(err).NotTo(HaveOccurred())

							job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
							Expect(err).NotTo(HaveOccurred())

							internalDomains, err := job.Property("app_domains")
							Expect(err).NotTo(HaveOccurred())

							Expect(internalDomains).To(Equal([]interface{}{
								"apps.example.com",
								"mesh.apps.example.com",
								[]interface{}{
									map[interface{}]interface{}{
										"name":     "apps.internal",
										"internal": true,
									},
								},
							}))

							temporaryIstioDomains, err := job.Property("copilot/temporary_istio_domains")
							Expect(err).NotTo(HaveOccurred())

							Expect(temporaryIstioDomains).To(Equal([]interface{}{
								"mesh.apps.example.com",
							}))
						})
					})

					Describe("configured", func() {
						var instanceGroup string
						BeforeEach(func() {
							if productName == "ert" {
								instanceGroup = "cloud_controller"
							} else {
								instanceGroup = "control"
							}
						})

						It("is properly set", func() {
							inputProperties := map[string]interface{}{
								".properties.istio":        "enable",
								".properties.istio_domain": "superspecial.istio.domain.com",
							}

							manifest, err := product.RenderManifest(inputProperties)
							Expect(err).NotTo(HaveOccurred())

							job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
							Expect(err).NotTo(HaveOccurred())

							internalDomains, err := job.Property("app_domains")
							Expect(err).NotTo(HaveOccurred())

							Expect(internalDomains).To(Equal([]interface{}{
								"apps.example.com",
								"superspecial.istio.domain.com",
								[]interface{}{
									map[interface{}]interface{}{
										"name":     "apps.internal",
										"internal": true,
									},
								},
							}))

							temporaryIstioDomains, err := job.Property("copilot/temporary_istio_domains")
							Expect(err).NotTo(HaveOccurred())

							Expect(temporaryIstioDomains).To(Equal([]interface{}{
								"superspecial.istio.domain.com",
							}))
						})
					})
				})
			})

			Context("when it is disabled", func() {
				It("zeros out istio-control, istio-router, and cc_route_syncer", func() {
					inputProperties := map[string]interface{}{}

					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					instanceCount, err := manifest.Path("/instance_groups/name=istio_control/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(0))

					instanceCount, err = manifest.Path("/instance_groups/name=istio_router/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(0))

					instanceCount, err = manifest.Path("/instance_groups/name=route_syncer/instances")
					Expect(err).NotTo(HaveOccurred())

					Expect(instanceCount).To(Equal(0))
				})

				It("Does not make an istio domain", func() {
					inputProperties := map[string]interface{}{}

					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(capiInstanceGroup, "cloud_controller_ng")
					Expect(err).NotTo(HaveOccurred())

					internalDomains, err := job.Property("app_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(internalDomains).To(Equal([]interface{}{
						"apps.example.com",
						[]interface{}{},
						[]interface{}{
							map[interface{}]interface{}{
								"name":     "apps.internal",
								"internal": true,
							},
						},
					}))

					temporaryIstioDomains, err := job.Property("copilot/temporary_istio_domains")
					Expect(err).NotTo(HaveOccurred())

					Expect(temporaryIstioDomains).To(Equal([]interface{}{
						[]interface{}{},
					}))
				})
			})
		})
	})

	Describe("Routing", func() {
		Describe("drain_timeout", func() {
			var (
				inputProperties     map[string]interface{}
				routerInstanceGroup string
			)

			BeforeEach(func() {
				routerInstanceGroup = "router"
			})

			Describe("when the property is set", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{
						".router.drain_timeout": 999,
					}
				})

				It("sets the prune_all_stale_routes", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
					Expect(err).NotTo(HaveOccurred())

					drainTimeout, err := job.Property("router/drain_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(drainTimeout).To(Equal(999))
				})
			})

			Describe("when the property is not set", func() {
				BeforeEach(func() {
					inputProperties = map[string]interface{}{}
				})

				It("defaults to false", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
					Expect(err).NotTo(HaveOccurred())

					drainTimeout, err := job.Property("router/drain_timeout")
					Expect(err).NotTo(HaveOccurred())
					Expect(drainTimeout).To(Equal(900))
				})
			})
		})

		Describe("gorouter", func() {
			var (
				inputProperties     map[string]interface{}
				routerInstanceGroup string
			)

			BeforeEach(func() {
				routerInstanceGroup = "router"
				inputProperties = map[string]interface{}{
					".properties.router_headers_remove_if_specified": []map[string]interface{}{
						{
							"name": "header1",
						},
						{
							"name": "header2",
						},
					}}
			})

			It("sets the headers to be removed for http responses", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(routerInstanceGroup, "gorouter")
				Expect(err).NotTo(HaveOccurred())

				removeHeaders, err := job.Property("router/http_rewrite/responses/remove_headers")
				Expect(err).NotTo(HaveOccurred())
				Expect(removeHeaders.([]interface{})[0].(map[interface{}]interface{})["name"]).To(Equal("header1"))
				Expect(removeHeaders.([]interface{})[1].(map[interface{}]interface{})["name"]).To(Equal("header2"))
			})
		})
	})
})
