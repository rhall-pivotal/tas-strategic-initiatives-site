package manifest_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("NFS volume service", func() {
	var instanceGroup string

	It("verifies that the nfsbrokerpush errand is configured correctly", func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}

		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		nfsBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
		Expect(err).NotTo(HaveOccurred())

		appDomain, err := nfsBrokerPush.Property("nfsbrokerpush/app_domain")
		Expect(err).NotTo(HaveOccurred())
		Expect(appDomain).To(Equal("sys.example.com"))

		nfsBrokerPushCredhubProperties, err := nfsBrokerPush.Property("nfsbrokerpush/credhub")
		Expect(err).NotTo(HaveOccurred())

		Expect(nfsBrokerPushCredhubProperties).To(HaveKeyWithValue("uaa_ca_cert", "fake-ops-manager-ca-certificate"))
		Expect(nfsBrokerPushCredhubProperties).To(HaveKeyWithValue("uaa_client_id", "nfs-broker-credhub"))
		Expect(nfsBrokerPushCredhubProperties).To(HaveKey("uaa_client_secret"))

		nfsBrokerPushDatabaseProperties, err := nfsBrokerPush.Property("nfsbrokerpush/db")
		Expect(err).NotTo(HaveOccurred())

		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("ca_cert", BeNil()))
		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("host", "mysql.service.cf.internal"))
		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("port", 3306))
		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("username", "((nfs-volume-db-credentials.username))"))
		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("password", "((nfs-volume-db-credentials.password))"))
		Expect(nfsBrokerPushDatabaseProperties).To(HaveKeyWithValue("name", "nfsvolume"))

		domain, err := nfsBrokerPush.Property("nfsbrokerpush/domain")
		Expect(err).NotTo(HaveOccurred())
		Expect(domain).To(Equal("sys.example.com"))

		org, err := nfsBrokerPush.Property("nfsbrokerpush/organization")
		Expect(err).NotTo(HaveOccurred())
		Expect(org).To(Equal("system"))

		password, err := nfsBrokerPush.Property("nfsbrokerpush/password")
		Expect(err).NotTo(HaveOccurred())
		Expect(password).To(Equal("((nfs-broker-push-db-credentials.password))"))

		skipCertVerify, err := nfsBrokerPush.Property("nfsbrokerpush/skip_cert_verify")
		Expect(err).NotTo(HaveOccurred())
		Expect(skipCertVerify).To(BeFalse())

		space, err := nfsBrokerPush.Property("nfsbrokerpush/space")
		Expect(err).NotTo(HaveOccurred())
		Expect(space).To(Equal("nfs"))

		storeID, err := nfsBrokerPush.Property("nfsbrokerpush/store_id")
		Expect(err).NotTo(HaveOccurred())
		Expect(storeID).To(Equal("nfsbroker"))

		syslogUrl, err := nfsBrokerPush.Property("nfsbrokerpush/syslog_url")
		Expect(err).NotTo(HaveOccurred())
		Expect(syslogUrl).To(BeEmpty())

		username, err := nfsBrokerPush.Property("nfsbrokerpush/username")
		Expect(err).NotTo(HaveOccurred())
		Expect(username).To(Equal("((nfs-broker-push-db-credentials.username))"))
	})

	It("creates a UAA client for the nfsbroker to use to access Credhub", func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "uaa"
		}

		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
		Expect(err).NotTo(HaveOccurred())

		nfsBrokerCredhubUaaClientProperties, err := uaa.Property("uaa/clients/nfs-broker-credhub")
		Expect(err).NotTo(HaveOccurred())

		Expect(nfsBrokerCredhubUaaClientProperties).To(HaveKeyWithValue("id", "nfs-broker-credhub"))
		Expect(nfsBrokerCredhubUaaClientProperties).To(HaveKey("authorities"))
		Expect(nfsBrokerCredhubUaaClientProperties).To(HaveKeyWithValue("authorized-grant-types", "client_credentials"))
		Expect(nfsBrokerCredhubUaaClientProperties).To(HaveKey("secret"))

		rawAuthorities, err := uaa.Property("uaa/clients/nfs-broker-credhub/authorities")
		Expect(err).ToNot(HaveOccurred())

		authorities := strings.Split(rawAuthorities.(string), ",")
		Expect(authorities).To(ConsistOf([]string{"credhub.read", "credhub.write"}))
	})

	It("grants permissions to the nfs-broker-credhub client", func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "credhub"
		}

		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
		Expect(err).NotTo(HaveOccurred())

		permissions, err := credhub.Property("credhub/authorization/permissions")
		Expect(err).ToNot(HaveOccurred())

		Expect(permissions).To(ContainElement(map[interface{}]interface{}{
			"path":       "/nfsbroker/*",
			"actors":     []interface{}{"uaa-client:nfs-broker-credhub"},
			"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
		}))
	})

	Context("when the syslog properties are configured", func() {
		var syslogConfiguration map[string]interface{}

		BeforeEach(func() {
			syslogConfiguration = map[string]interface{}{
				".properties.syslog_host":     "example.com",
				".properties.syslog_port":     12345,
				".properties.syslog_protocol": "tcp",
			}
		})

		It("configures the syslog_url property", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderManifest(syslogConfiguration)
			Expect(err).NotTo(HaveOccurred())

			nfsBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			syslogUrl, err := nfsBrokerPush.Property("nfsbrokerpush/syslog_url")
			Expect(err).NotTo(HaveOccurred())
			Expect(syslogUrl).To(Equal("tcp://example.com:12345"))
		})
	})

	Context("when the NFS V3 driver is enabled without LDAP configuration", func() {
		It("disables LDAP on the nfsbrokerpush job", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			nfsBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			ldapEnabled, err := nfsBrokerPush.Property("nfsbrokerpush/ldap_enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapEnabled).To(BeFalse())
		})

		It("does not configure LDAP on the nfsv3driver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			nfsV3Driver, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverProperties, err := nfsV3Driver.Property("nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("disable", BeFalse()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_svc_user", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_svc_password", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_host", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_port", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_user_fqdn", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_ca_cert", BeNil()))
			Expect(nfsV3DriverProperties).NotTo(HaveKey("allowed-in-source"))
		})
	})

	Context("when the NFS V3 driver is enabled with LDAP configuration", func() {
		var ldapConfiguration map[string]interface{}

		BeforeEach(func() {
			ldapConfiguration = map[string]interface{}{
				".properties.nfs_volume_driver.enable.ldap_service_account_user": "service-account-user",
				".properties.nfs_volume_driver.enable.ldap_service_account_password": map[string]string{
					"secret": "service-account-password",
				},
				".properties.nfs_volume_driver.enable.ldap_server_host": "ldap-host",
				".properties.nfs_volume_driver.enable.ldap_server_port": 12345,
				".properties.nfs_volume_driver.enable.ldap_user_fqdn":   "ldap-user-search-base",
			}
		})

		It("enables LDAP on the nfsbrokerpush job", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderManifest(ldapConfiguration)
			Expect(err).NotTo(HaveOccurred())

			nfsBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			ldapEnabled, err := nfsBrokerPush.Property("nfsbrokerpush/ldap_enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapEnabled).To(BeTrue())
		})

		It("configures LDAP on the nfsv3driver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderManifest(ldapConfiguration)
			Expect(err).NotTo(HaveOccurred())

			nfsV3Driver, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverProperties, err := nfsV3Driver.Property("nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("disable", BeFalse()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_svc_user", "service-account-user"))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_svc_password", Not(BeEmpty())))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_host", "ldap-host"))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_port", 12345))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_user_fqdn", "ldap-user-search-base"))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_ca_cert", BeNil()))
			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("allowed-in-source", ""))
		})

		Context("when the LDAP CA certificate is configured", func() {
			It("configures LDAP on the nfsv3driver job", func() {
				if productName == "srt" {
					instanceGroup = "compute"
				} else {
					instanceGroup = "diego_cell"
				}

				ldapConfiguration[".properties.nfs_volume_driver.enable.ldap_ca_cert"] = "ldap-ca-cert"

				manifest, err := product.RenderManifest(ldapConfiguration)
				Expect(err).NotTo(HaveOccurred())

				nfsV3Driver, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
				Expect(err).NotTo(HaveOccurred())

				nfsV3DriverProperties, err := nfsV3Driver.Property("nfsv3driver")
				Expect(err).NotTo(HaveOccurred())

				Expect(nfsV3DriverProperties).To(HaveKeyWithValue("ldap_ca_cert", "ldap-ca-cert"))
			})
		})
	})

	Context("when the NFS V3 driver is disabled", func() {
		It("disables the nfsv3driver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.nfs_volume_driver": "disable",
			})
			Expect(err).NotTo(HaveOccurred())

			nfsV3Driver, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverProperties, err := nfsV3Driver.Property("nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			Expect(nfsV3DriverProperties).To(HaveKeyWithValue("disable", BeTrue()))

			mapfsJob, err := manifest.FindInstanceGroupJob(instanceGroup, "mapfs")
			Expect(err).NotTo(HaveOccurred())

			mapfsDisableProperty, err := mapfsJob.Property("disable")
			Expect(err).NotTo(HaveOccurred())

			Expect(mapfsDisableProperty).To(BeTrue())
		})
	})

	Context("when the NFS volumes services are enabled", func() {
		It("configures the nfsbrokerpush job", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.nfs_volume_driver": "enable",
			})
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			testNfsBrokerPushProperties(job)
			Expect(job.Path("/provides/nfsbrokerpush")).To(Equal(map[interface{}]interface{}{"as": "ignore-me"}))
		})
	})

	Describe("Backup and Restore", func() {
		Context("on the backup_restore instance group", func() {
			instanceGroup := "backup_restore"

			It("configures the nfsbrokerpush job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
				Expect(err).NotTo(HaveOccurred())

				testNfsBrokerPushProperties(job)

				nfsBrokerPushCredhubProperties, err := job.Property("nfsbrokerpush/credhub")
				Expect(err).NotTo(HaveOccurred())

				Expect(nfsBrokerPushCredhubProperties).To(HaveKeyWithValue("uaa_ca_cert", "fake-ops-manager-ca-certificate"))
				Expect(nfsBrokerPushCredhubProperties).To(HaveKeyWithValue("uaa_client_id", "nfs-broker-credhub"))
				Expect(nfsBrokerPushCredhubProperties).To(HaveKey("uaa_client_secret"))

				storeID, err := job.Property("nfsbrokerpush/store_id")
				Expect(err).NotTo(HaveOccurred())
				Expect(storeID).To(Equal("nfsbroker"))

				Expect(job.Path("/provides/nfsbrokerpush")).To(Equal(map[interface{}]interface{}{"as": "nfsbrokerpush"}))
			})

			It("configures the nfsbroker-bbr-lock job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbroker-bbr-lock")
				Expect(err).NotTo(HaveOccurred())
				Expect(job.Path("/consumes/nfsbrokerpush")).To(Equal(map[interface{}]interface{}{"from": "nfsbrokerpush"}))
			})
		})
	})
})

func testNfsBrokerPushProperties(nfsBrokerPush planitest.Manifest) {
	nfsBrokerPushProperties, err := nfsBrokerPush.Path("/properties/nfsbrokerpush")
	Expect(err).NotTo(HaveOccurred())
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("domain", "sys.example.com"))
	Expect(nfsBrokerPushProperties).To(HaveKey("db"))
	Expect(nfsBrokerPushProperties).To(HaveKey("cf"))
	nfsBrokerCFProperties := (nfsBrokerPushProperties.(map[interface{}]interface{}))["cf"]
	Expect(nfsBrokerCFProperties).To(HaveKeyWithValue("admin_user", "admin"))
	Expect(nfsBrokerCFProperties).To(HaveKey("admin_password"))
	Expect(nfsBrokerCFProperties).To(HaveKey("dial_timeout"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("organization", "system"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("skip_cert_verify", BeFalse()))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("app_domain", "sys.example.com"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("space", "nfs"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("username", "((nfs-broker-push-db-credentials.username))"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("password", "((nfs-broker-push-db-credentials.password))"))
	Expect(nfsBrokerPushProperties).To(HaveKeyWithValue("syslog_url", ""))
	Expect(nfsBrokerPushProperties).To(HaveKey("ldap_enabled"))
}
