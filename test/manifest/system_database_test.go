package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("System Database", func() {
	Describe("External Database", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "diego_database"
			} else {
				instanceGroup = "control"
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

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "policy-server")
			Expect(err).NotTo(HaveOccurred())

			requireSSL, err := job.Property("database/require_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(requireSSL).To(BeFalse())

			dbType, err := job.Property("database/type")
			Expect(err).NotTo(HaveOccurred())
			Expect(dbType).To(Equal("mysql"))

			host, err := job.Property("database/host")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("foo.bar"))

			port, err := job.Property("database/port")
			Expect(err).NotTo(HaveOccurred())
			Expect(port).To(Equal(5432))
		})

		Context("when the operator provides a CA certificate", func() {
			BeforeEach(func() {
				inputProperties[".properties.system_database.external.ca_cert"] = "fake-ca-cert"
			})

			It("configures jobs to use that CA certificate ", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				jobs := []string{"policy-server", "silk-controller"}

				for _, j := range jobs {
					job, err := manifest.FindInstanceGroupJob(instanceGroup, j)
					Expect(err).NotTo(HaveOccurred())

					requireSSL, err := job.Property("database/require_ssl")
					Expect(err).NotTo(HaveOccurred())
					Expect(requireSSL).To(BeTrue())

					caCert, err := job.Property("database/ca_cert")
					Expect(err).NotTo(HaveOccurred())
					Expect(caCert).To(Equal("fake-ca-cert"))
				}

				jobs = []string{"locket", "bbs"}

				for _, j := range jobs {
					job, err := manifest.FindInstanceGroupJob(instanceGroup, j)
					Expect(err).NotTo(HaveOccurred())

					requireSSL, err := job.Property(fmt.Sprintf("diego/%s/sql/require_ssl", j))
					Expect(err).NotTo(HaveOccurred())
					Expect(requireSSL).To(BeTrue())

					caCert, err := job.Property(fmt.Sprintf("diego/%s/sql/ca_cert", j))
					Expect(err).NotTo(HaveOccurred())
					Expect(caCert).To(Equal("fake-ca-cert"))
				}
			})
		})
	})
})
