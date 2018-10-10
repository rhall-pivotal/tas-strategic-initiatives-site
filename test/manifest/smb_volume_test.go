package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SMB volume service", func() {
	var instanceGroup string

	Context("when the SMB driver is enabled", func() {
		It("enables the smbdriver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			smbDriver, err := manifest.FindInstanceGroupJob(instanceGroup, "smbdriver")
			Expect(err).NotTo(HaveOccurred())

			smbDriverDisabled, err := smbDriver.Property("disable")
			Expect(err).NotTo(HaveOccurred())

			Expect(smbDriverDisabled).To(BeFalse())
		})

		It("configures new UAA clients", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "uaa"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			uaaClients, err := uaa.Property("uaa/clients")
			Expect(err).NotTo(HaveOccurred())

			Expect(uaaClients).To(HaveKey("smb-broker"))
			uaaSmbBrokerClient := (uaaClients.(map[interface{}]interface{}))["smb-broker"]
			Expect(uaaSmbBrokerClient).To(HaveKeyWithValue("id", "smb-broker"))
			Expect(uaaSmbBrokerClient).To(HaveKeyWithValue("authorities", "cloud_controller.admin"))
			Expect(uaaSmbBrokerClient).To(HaveKeyWithValue("authorized-grant-types", "client_credentials"))
			Expect(uaaSmbBrokerClient).To(HaveKey("secret"))

			Expect(uaaClients).To(HaveKey("smb-broker-credhub"))
			uaaSmbBrokerCredhubClient := (uaaClients.(map[interface{}]interface{}))["smb-broker-credhub"]
			Expect(uaaSmbBrokerCredhubClient).To(HaveKeyWithValue("id", "smb-broker-credhub"))
			Expect(uaaSmbBrokerCredhubClient).To(HaveKeyWithValue("authorities", "credhub.read,credhub.write"))
			Expect(uaaSmbBrokerCredhubClient).To(HaveKeyWithValue("authorized-grant-types", "client_credentials"))
			Expect(uaaSmbBrokerCredhubClient).To(HaveKey("secret"))
		})

		It("configures credhub to allow access to the stored credentials", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "credhub"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			credhub, err := manifest.FindInstanceGroupJob(instanceGroup, "credhub")
			Expect(err).NotTo(HaveOccurred())

			credhubPermissions, err := credhub.Property("credhub/authorization/permissions")
			Expect(err).NotTo(HaveOccurred())

			Expect(credhubPermissions).To(ContainElement(map[interface{}]interface{}{
				"actors":     []interface{}{"uaa-client:smb-broker-credhub"},
				"operations": []interface{}{"read", "write", "delete", "read_acl", "write_acl"},
				"path":       "/smbbroker/*",
			}))
		})

		It("configures the smbbrokerpush errand", func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "clock_global"
			}

			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			smbBrokerPush, err := manifest.FindInstanceGroupJob(instanceGroup, "smbbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			smbBrokerProperties, err := smbBrokerPush.Path("/properties")
			Expect(err).NotTo(HaveOccurred())

			Expect(smbBrokerProperties).To(HaveKeyWithValue("app_domain", "sys.example.com"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("domain", "sys.example.com"))
			Expect(smbBrokerProperties).To(HaveKey("cf"))
			smbBrokerCFProperties := (smbBrokerProperties.(map[interface{}]interface{}))["cf"]
			Expect(smbBrokerCFProperties).To(HaveKeyWithValue("client_id", "smb-broker"))
			Expect(smbBrokerCFProperties).To(HaveKey("client_secret"))
			Expect(smbBrokerProperties).To(HaveKey("credhub"))
			smbBrokerCredhubProperties := (smbBrokerProperties.(map[interface{}]interface{}))["credhub"]
			Expect(smbBrokerCredhubProperties).To(HaveKeyWithValue("url", "https://credhub.service.cf.internal:8844"))
			Expect(smbBrokerCredhubProperties).To(HaveKeyWithValue("uaa_client_id", "smb-broker-credhub"))
			Expect(smbBrokerCredhubProperties).To(HaveKey("uaa_client_secret"))
			Expect(smbBrokerCredhubProperties).To(HaveKeyWithValue("store_id", "smbbroker"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("organization", "system"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("space", "smb"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("username", "((smb-broker-credentials.username))"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("password", "((smb-broker-credentials.password))"))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("skip_cert_verify", BeFalse()))
			Expect(smbBrokerProperties).To(HaveKeyWithValue("syslog_url", ""))
		})
	})

	Context("when the SMB driver is disabled", func() {
		It("disables the smbdriver job", func() {
			if productName == "srt" {
				instanceGroup = "compute"
			} else {
				instanceGroup = "diego_cell"
			}

			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.smb_volume_driver": "disable",
			})
			Expect(err).NotTo(HaveOccurred())

			smbDriver, err := manifest.FindInstanceGroupJob(instanceGroup, "smbdriver")
			Expect(err).NotTo(HaveOccurred())

			smbDriverDisabled, err := smbDriver.Property("disable")
			Expect(err).NotTo(HaveOccurred())

			Expect(smbDriverDisabled).To(BeTrue())
		})
	})
})
