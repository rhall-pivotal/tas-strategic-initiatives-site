package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("MySQL", func() {
	var (
		instanceGroup string
	)

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "database"
		} else {
			instanceGroup = "mysql"
		}
	})

	Context("when the operator selects clustered mysql", func() {
		var (
			inputProperties map[string]interface{}
			manifest        planitest.Manifest
		)

		BeforeEach(func() {
			inputProperties = map[string]interface{}{
				".properties.system_database": "internal_pxc",
				".mysql.max_connections":      40000,
			}
			var err error
			manifest, err = product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())
		})

		It("configures max connections for pxc-mysql to be the configured value", func() {
			mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
			Expect(err).NotTo(HaveOccurred())

			maxConnections, err := mysqlClustered.Property("engine_config/max_connections")
			Expect(err).NotTo(HaveOccurred())
			Expect(maxConnections).To(Equal(40000))
		})

		It("configures the port", func() {
			mysql, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
			Expect(err).NotTo(HaveOccurred())

			port, err := mysql.Property("port")
			Expect(err).NotTo(HaveOccurred())

			canary, err := manifest.FindInstanceGroupJob("mysql_monitor", "replication-canary")
			Expect(err).NotTo(HaveOccurred())

			canaryPort, err := canary.Property("mysql-monitoring/replication-canary/mysql_port")
			Expect(err).NotTo(HaveOccurred())

			if productName == "srt" {
				Expect(port).To(Equal(13306))
				Expect(canaryPort).To(Equal(13306))
			} else {
				Expect(port).To(Equal(3306))
				Expect(canaryPort).To(Equal(3306))
			}
		})
	})
})
