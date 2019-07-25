package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("MySQL Proxy", func() {
	var (
		instanceGroup string
	)
	Context("when the operator selects mysql proxy", func() {
		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "database"
			} else {
				instanceGroup = "mysql_proxy"
			}
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

			It("configures the max open files for the proxy", func() {
				proxyManifest, err := manifest.FindInstanceGroupJob(instanceGroup, "proxy")
				Expect(err).NotTo(HaveOccurred())

				maxConnections, err := proxyManifest.Property("max_open_files")
				Expect(err).NotTo(HaveOccurred())
				Expect(maxConnections).To(Equal(1048576))
			})
		})

		Context("when the operator chooses to expose the inactive node port", func() {
			var (
				inputProperties map[string]interface{}
				manifest        planitest.Manifest
			)

			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_database":             "internal_pxc",
					".mysql_proxy.enable_inactive_mysql_port": true,
				}
				var err error
				manifest, err = product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())
			})

			It("enables the inactive node port on the proxy", func() {
				proxyManifest, err := manifest.FindInstanceGroupJob(instanceGroup, "proxy")
				Expect(err).NotTo(HaveOccurred())

				inactiveMySQLPort, err := proxyManifest.Property("inactive_mysql_port")
				Expect(err).NotTo(HaveOccurred())
				Expect(inactiveMySQLPort).To(Equal(3336))
			})
		})

		Context("when the operator does not expose the inactive node port", func() {
			var (
				inputProperties map[string]interface{}
				manifest        planitest.Manifest
			)

			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_database":             "internal_pxc",
					".mysql_proxy.enable_inactive_mysql_port": false,
				}
				var err error
				manifest, err = product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())
			})

			It("leaves the property empty", func() {
				proxyManifest, err := manifest.FindInstanceGroupJob(instanceGroup, "proxy")
				Expect(err).NotTo(HaveOccurred())

				inactiveMySQLPort, err := proxyManifest.Property("inactive_mysql_port")
				Expect(err).NotTo(HaveOccurred())
				Expect(inactiveMySQLPort).To(BeNil())
			})
		})
	})
})
