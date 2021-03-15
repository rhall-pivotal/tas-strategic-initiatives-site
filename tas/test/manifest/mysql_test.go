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
	Context("when the operator selects clustered mysql", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "database"
			} else {
				instanceGroup = "mysql"
			}
		})

		Context("and configures the values for the properties", func() {
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

			It("configures origin tag for loggregator_agent", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "loggregator_agent")
				Expect(err).NotTo(HaveOccurred())

				tags, err := mysqlClustered.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("origin", "mysql"))
			})

			It("configures origin tag for loggr-forwarder-agent", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "loggr-forwarder-agent")
				Expect(err).NotTo(HaveOccurred())

				tags, err := mysqlClustered.Property("tags")
				Expect(err).NotTo(HaveOccurred())
				Expect(tags).To(HaveKeyWithValue("origin", "mysql"))
			})
		})
		Context("and uses the defaults", func() {
			var (
				inputProperties map[string]interface{}
				manifest        planitest.Manifest
			)

			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_database": "internal_pxc",
				}
				var err error
				manifest, err = product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())
			})

			It("sets the max connections for pxc-mysql to be the configured value", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
				Expect(err).NotTo(HaveOccurred())

				maxConnections, err := mysqlClustered.Property("engine_config/max_connections")
				Expect(err).NotTo(HaveOccurred())
				Expect(maxConnections).To(Equal(3500))
			})

			It("sets the max open files for the mysql", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
				Expect(err).NotTo(HaveOccurred())

				maxConnections, err := mysqlClustered.Property("max_open_files")
				Expect(err).NotTo(HaveOccurred())
				Expect(maxConnections).To(Equal(1048576))
			})

			It("sets the innodb buffer pool size to 50% of available memory", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
				Expect(err).NotTo(HaveOccurred())

				maxConnections, err := mysqlClustered.Property("engine_config/innodb_buffer_pool_size_percent")
				Expect(err).NotTo(HaveOccurred())
				Expect(maxConnections).To(Equal(50))
			})
		})

		Context("logging_timestamp_format", func() {
			var proxyInstanceGroup string

			BeforeEach(func() {
				if productName == "srt" {
					proxyInstanceGroup = "database"
				} else {
					proxyInstanceGroup = "mysql_proxy"
				}
			})

			When("logging_timestamp_format is set to deprecated", func() {
				It("is used in replication-canary and mysql-metrics", func() {
					manifest := renderProductManifest(product, map[string]interface{}{
						".properties.logging_timestamp_format": "deprecated",
					})

					replicationCanary, err := manifest.FindInstanceGroupJob("mysql_monitor", "replication-canary")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err := replicationCanary.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))

					mysqlMetrics, err := manifest.FindInstanceGroupJob(instanceGroup, "mysql-metrics")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = mysqlMetrics.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))

					galeraAgent, err := manifest.FindInstanceGroupJob(instanceGroup, "galera-agent")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = galeraAgent.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))

					graLogPurger, err := manifest.FindInstanceGroupJob(instanceGroup, "gra-log-purger")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = graLogPurger.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))

					proxy, err := manifest.FindInstanceGroupJob(proxyInstanceGroup, "proxy")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = proxy.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))

					pxcMysql, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = pxcMysql.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("unix-epoch"))
				})
			})

			When("logging_format_timestamp is set to rfc3339", func() {
				It("is used in replication-canary and mysql-metrics", func() {
					manifest := renderProductManifest(product, map[string]interface{}{
						".properties.logging_timestamp_format": "rfc3339",
					})

					replicationCanary, err := manifest.FindInstanceGroupJob("mysql_monitor", "replication-canary")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err := replicationCanary.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))

					mysqlMetrics, err := manifest.FindInstanceGroupJob(instanceGroup, "mysql-metrics")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = mysqlMetrics.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))

					galeraAgent, err := manifest.FindInstanceGroupJob(instanceGroup, "galera-agent")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = galeraAgent.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))

					graLogPurger, err := manifest.FindInstanceGroupJob(instanceGroup, "gra-log-purger")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = graLogPurger.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))

					proxy, err := manifest.FindInstanceGroupJob(proxyInstanceGroup, "proxy")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = proxy.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))

					pxcMysql, err := manifest.FindInstanceGroupJob(instanceGroup, "pxc-mysql")
					Expect(err).NotTo(HaveOccurred())
					loggingFormatTimestamp, err = pxcMysql.Property("logging/format/timestamp")
					Expect(err).NotTo(HaveOccurred())
					Expect(loggingFormatTimestamp).To(Equal("rfc3339"))
				})
			})
		})
	})
})
