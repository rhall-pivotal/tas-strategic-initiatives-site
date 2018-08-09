package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NFS volume service", func() {
	var instanceGroup string

	Context("when the NFS V3 driver is enabled without LDAP configuration", func() {
		It("enables the nfsv3driver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverDisable, err := nfsV3DriverPush.Property("nfsv3driver/disable")
			Expect(err).NotTo(HaveOccurred())
			Expect(nfsV3DriverDisable).To(BeFalse())
		})

		It("disables LDAP on the nfsbrokerpush job", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			nfsBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			ldapEnabled, err := nfsBrokerPush.Property("nfsbrokerpush/ldap_enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapEnabled).To(BeFalse())
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

			manifest, err := product.RenderService.RenderManifest(ldapConfiguration)
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

			manifest, err := product.RenderService.RenderManifest(ldapConfiguration)
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			ldapServiceUser, err := nfsV3DriverPush.Property("nfsv3driver/ldap_svc_user")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapServiceUser).To(Equal("service-account-user"))

			ldapServicePassword, err := nfsV3DriverPush.Property("nfsv3driver/ldap_svc_password")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapServicePassword).To(MatchRegexp("((/opsmgr/cf-[a-z0-9]{20}/nfs_volume_driver/enable/ldap_service_account_password.value))"))

			ldapHost, err := nfsV3DriverPush.Property("nfsv3driver/ldap_host")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapHost).To(Equal("ldap-host"))

			ldapPort, err := nfsV3DriverPush.Property("nfsv3driver/ldap_port")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapPort).To(Equal(12345))

			ldapUserFqdn, err := nfsV3DriverPush.Property("nfsv3driver/ldap_user_fqdn")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapUserFqdn).To(Equal("ldap-user-search-base"))

			ldapCACert, err := nfsV3DriverPush.Property("nfsv3driver/ldap_ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(ldapCACert).To(BeNil())
		})

		Context("when the LDAP CA certificate is configured", func() {
			It("configures LDAP on the nfsv3driver job", func() {
				if productName == "srt" {
					instanceGroup = "compute"
				} else {
					instanceGroup = "diego_cell"
				}

				ldapConfiguration[".properties.nfs_volume_driver.enable.ldap_ca_cert"] = "ldap-ca-cert"

				manifest, err := product.RenderService.RenderManifest(ldapConfiguration)
				Expect(err).NotTo(HaveOccurred())

				nfsV3DriverPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
				Expect(err).NotTo(HaveOccurred())

				ldapCACert, err := nfsV3DriverPush.Property("nfsv3driver/ldap_ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(ldapCACert).To(Equal("ldap-ca-cert"))
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

			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
				".properties.nfs_volume_driver": "disable",
			})
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverPush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsv3driver")
			Expect(err).NotTo(HaveOccurred())

			nfsV3DriverDisable, err := nfsV3DriverPush.Property("nfsv3driver/disable")
			Expect(err).NotTo(HaveOccurred())
			Expect(nfsV3DriverDisable).To(BeTrue())
		})
	})
})
