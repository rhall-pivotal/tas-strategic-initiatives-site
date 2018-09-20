package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("System Database", func() {
	Describe("External Database", func() {
		var (
			inputProperties         map[string]interface{}
			controllerInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "ert" {
				controllerInstanceGroup = "diego_database"
			} else {
				controllerInstanceGroup = "control"
			}
			inputProperties = map[string]interface{}{
				".properties.system_database":                                       "external",
				".properties.system_database.external.host":                         "foo.bar",
				".properties.system_database.external.port":                         5432,
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

		It("configures jobs with user provided values", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(controllerInstanceGroup, "policy-server")
			Expect(err).NotTo(HaveOccurred())

			property, err := job.Property("database/type")
			Expect(err).NotTo(HaveOccurred())
			Expect(property).To(Equal("mysql"))

			property, err = job.Property("database/host")
			Expect(err).NotTo(HaveOccurred())
			Expect(property).To(Equal("foo.bar"))

			property, err = job.Property("database/port")
			Expect(err).NotTo(HaveOccurred())
			Expect(property).To(Equal(5432))
		})
	})
})
