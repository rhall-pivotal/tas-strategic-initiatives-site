package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("System Database", func() {
	Describe("Internal PXC", func() {
		var (
			inputProperties map[string]interface{}
			instanceGroup   string
		)

		BeforeEach(func() {
			if productName == "ert" {
				instanceGroup = "clock_global"
			} else {
				instanceGroup = "control"
			}
			inputProperties = map[string]interface{}{}
		})

		It("configures the jobs", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			// nfsbroker
			nfsbrokerpush, err := manifest.FindInstanceGroupJob(instanceGroup, "nfsbrokerpush")
			Expect(err).NotTo(HaveOccurred())

			host, err := nfsbrokerpush.Property("nfsbrokerpush/db/host")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("mysql.service.cf.internal"))

			caCert, err := nfsbrokerpush.Property("nfsbrokerpush/db/ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).ToNot(BeEmpty())

			nfsbrokerbbr, err := manifest.FindInstanceGroupJob("backup_restore", "nfsbroker-bbr")
			Expect(err).NotTo(HaveOccurred())

			host, err = nfsbrokerbbr.Property("nfsbroker/db_hostname")
			Expect(err).NotTo(HaveOccurred())
			Expect(host).To(Equal("mysql.service.cf.internal"))

			backup, err := nfsbrokerbbr.Property("nfsbroker/release_level_backup")
			Expect(err).NotTo(HaveOccurred())
			Expect(backup).To(BeTrue())

			caCert, err = nfsbrokerbbr.Property("nfsbroker/db_ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).ToNot(BeEmpty())

			// usage-service
			pushUsageService, err := manifest.FindInstanceGroupJob(instanceGroup, "push-usage-service")
			Expect(err).NotTo(HaveOccurred())

			caCert, err = pushUsageService.Property("databases/app_usage_service/ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).NotTo(BeEmpty())

			verifySSL, err := pushUsageService.Property("databases/app_usage_service/verify_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(verifySSL).To(BeTrue())

			bbrUsageServiceDB, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-usage-servicedb")
			Expect(err).NotTo(HaveOccurred())

			caCert, err = bbrUsageServiceDB.Property("database/ca_cert")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCert).NotTo(BeEmpty())

			verifySSL, err = bbrUsageServiceDB.Property("database/skip_host_verify")
			Expect(err).NotTo(HaveOccurred())
			Expect(verifySSL).To(BeTrue())
		})
	})

	Describe("External Database", func() {
		var (
			inputProperties map[string]interface{}
			dbInstanceGroup string
			ccInstanceGroup string
			cgInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "ert" {
				dbInstanceGroup = "diego_database"
				ccInstanceGroup = "cloud_controller"
				cgInstanceGroup = "clock_global"
			} else {
				dbInstanceGroup = "control"
				ccInstanceGroup = "control"
				cgInstanceGroup = "control"
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

			job, err := manifest.FindInstanceGroupJob(dbInstanceGroup, "policy-server")
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

			// usage-service should not verify SSL
			pushUsageService, err := manifest.FindInstanceGroupJob(cgInstanceGroup, "push-usage-service")
			Expect(err).NotTo(HaveOccurred())

			verifySSL, err := pushUsageService.Property("databases/app_usage_service/verify_ssl")
			Expect(err).NotTo(HaveOccurred())
			Expect(verifySSL).To(BeFalse())

			bbrUsageServiceDB, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-usage-servicedb")
			Expect(err).NotTo(HaveOccurred())

			skipHostVerify, err := bbrUsageServiceDB.Property("database/skip_host_verify")
			Expect(err).NotTo(HaveOccurred())
			Expect(skipHostVerify).To(BeFalse())
		})

		Context("when the operator provides a CA certificate", func() {
			BeforeEach(func() {
				inputProperties[".properties.system_database.external.ca_cert"] = "fake-ca-cert"
			})

			It("configures jobs to use that CA certificate ", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				// policy-server
				job, err := manifest.FindInstanceGroupJob(dbInstanceGroup, "policy-server")
				Expect(err).NotTo(HaveOccurred())

				requireSSL, err := job.Property("database/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireSSL).To(BeTrue())

				caCert, err := job.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// silk-controller
				job, err = manifest.FindInstanceGroupJob(dbInstanceGroup, "silk-controller")
				Expect(err).NotTo(HaveOccurred())

				requireSSL, err = job.Property("database/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireSSL).To(BeTrue())

				caCert, err = job.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// locket
				job, err = manifest.FindInstanceGroupJob(dbInstanceGroup, "locket")
				Expect(err).NotTo(HaveOccurred())

				requireSSL, err = job.Property("diego/locket/sql/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireSSL).To(BeTrue())

				caCert, err = job.Property("diego/locket/sql/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// bbs
				job, err = manifest.FindInstanceGroupJob(dbInstanceGroup, "bbs")
				Expect(err).NotTo(HaveOccurred())

				requireSSL, err = job.Property("diego/bbs/sql/require_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireSSL).To(BeTrue())

				caCert, err = job.Property("diego/bbs/sql/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// cloud_controller_ng
				job, err = manifest.FindInstanceGroupJob(ccInstanceGroup, "cloud_controller_ng")
				Expect(err).NotTo(HaveOccurred())

				requireSSL, err = job.Property("ccdb/ssl_verify_hostname")
				Expect(err).NotTo(HaveOccurred())
				Expect(requireSSL).To(BeTrue())

				caCert, err = job.Property("ccdb/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// routing-api
				job, err = manifest.FindInstanceGroupJob(ccInstanceGroup, "routing-api")
				Expect(err).NotTo(HaveOccurred())

				caCert, err = job.Property("routing_api/sqldb/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// nfsbroker
				nfsbrokerpush, err := manifest.FindInstanceGroupJob(cgInstanceGroup, "nfsbrokerpush")
				Expect(err).NotTo(HaveOccurred())

				caCert, err = nfsbrokerpush.Property("nfsbrokerpush/db/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				nfsbrokerbbr, err := manifest.FindInstanceGroupJob("backup_restore", "nfsbroker-bbr")
				Expect(err).NotTo(HaveOccurred())

				caCert, err = nfsbrokerbbr.Property("nfsbroker/db_ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				// usage-service
				pushUsageService, err := manifest.FindInstanceGroupJob(cgInstanceGroup, "push-usage-service")
				Expect(err).NotTo(HaveOccurred())

				caCert, err = pushUsageService.Property("databases/app_usage_service/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				verifySSL, err := pushUsageService.Property("databases/app_usage_service/verify_ssl")
				Expect(err).NotTo(HaveOccurred())
				Expect(verifySSL).To(BeTrue())

				bbrUsageServiceDB, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-usage-servicedb")
				Expect(err).NotTo(HaveOccurred())

				caCert, err = bbrUsageServiceDB.Property("database/ca_cert")
				Expect(err).NotTo(HaveOccurred())
				Expect(caCert).To(Equal("fake-ca-cert"))

				skipHostVerify, err := bbrUsageServiceDB.Property("database/skip_host_verify")
				Expect(err).NotTo(HaveOccurred())
				Expect(skipHostVerify).To(BeTrue())
			})
		})
	})
})
