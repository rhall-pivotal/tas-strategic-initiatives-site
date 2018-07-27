package manifest_test

import (
	"github.com/pivotal-cf/planitest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

	Describe("PXC releases", func() {
		It("colocates the cluster-health-logger job on the appropriate instance", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "cluster-health-logger")
			Expect(err).NotTo(HaveOccurred())
		})

		It("colocates the galera-agent job on the appropriate instance", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "galera-agent")
			Expect(err).NotTo(HaveOccurred())
		})

		It("colocates the gra-log-purger job on the appropriate instance", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "gra-log-purger")
			Expect(err).NotTo(HaveOccurred())
		})

		It("colocates the pxc-mysql job on the appropriate instance", func() {
			manifest, err := product.RenderService.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
			Expect(err).NotTo(HaveOccurred())
		})
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

	Describe("when the operator configures max connections for mysql", func() {
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
			BeforeEach(func() {
				manifest, err = product.RenderService.RenderManifest(map[string]interface{}{
					".properties.system_database": "internal_pxc",
					".mysql.max_connections":      40000,
				})
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
