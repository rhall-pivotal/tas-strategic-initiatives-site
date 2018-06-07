package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/planitest"
)

var _ = Describe("MySQL", func() {
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
			mysql, err := manifest.FindInstanceGroupJob("mysql", "mysql")
			Expect(err).NotTo(HaveOccurred())
			Expect(mysql.Property("cf_mysql/mysql/max_connections")).To(Equal(10000))
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

			It("configures max connections for mysql-clustered to be the configured value", func() {
				mysqlClustered, err := manifest.FindInstanceGroupJob("mysql", "mysql-clustered")
				Expect(err).NotTo(HaveOccurred())
				Expect(mysqlClustered.Property("max_connections")).To(Equal(40000))
			})
		})
	})
})
