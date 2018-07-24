package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("MySQL", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "database"
		} else {
			instanceGroup = "mysql"
		}
	})

	Describe("when the operator turns on audit logging", func() {
		It("enables audit logs", func() {
			manifest, err := product.RenderService.RenderManifest(map[string]interface{}{
				".properties.mysql_activity_logging": "enable",
				".properties.system_database":        "internal_pxc",
			})
			Expect(err).NotTo(HaveOccurred())

			mysql, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
			Expect(err).NotTo(HaveOccurred())

			auditLogsEnabled, err := mysql.Property("engine_config/audit_logs/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(auditLogsEnabled).To(BeTrue())
		})
	})

	Context("when the operator configures max connections for mysql", func() {
		var (
			manifest planitest.Manifest
			err      error
		)

		BeforeEach(func() {
			manifest, err = product.RenderService.RenderManifest(map[string]interface{}{
				".mysql.max_connections": 10000,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("configures the max connections for mysql to be the set value", func() {
			mysql, err := manifest.FindInstanceGroupJob(instanceGroup, "mysql")
			Expect(err).NotTo(HaveOccurred())

			maxConnections, err := mysql.Property("cf_mysql/mysql/max_connections")
			Expect(err).NotTo(HaveOccurred())
			Expect(maxConnections).To(Equal(10000))
		})

		Context("when the operator selects clustered mysql", func() {
			var (
				inputProperties map[string]interface{}
			)

			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_database": "internal_pxc",
					".mysql.max_connections":      40000,
				}
				manifest, err = product.RenderService.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())
			})

			It("configures max connections for pxc-mysql to be the configured value", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
				Expect(err).NotTo(HaveOccurred())

				maxConnections, err := mysqlClustered.Property("engine_config/max_connections")
				Expect(err).NotTo(HaveOccurred())
				Expect(maxConnections).To(Equal(40000))
			})
		})
	})
})
