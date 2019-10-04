package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("winc_network", func() {
	It("sets networking properties", func() {
		manifest, err := product.RenderManifest(map[string]interface{}{})
		Expect(err).NotTo(HaveOccurred())

		winc, err := manifest.FindInstanceGroupJob("windows_diego_cell", "winc-network-hns-acls")
		Expect(err).NotTo(HaveOccurred())

		dnsServers, err := winc.Property("winc_network/dns_servers")
		Expect(err).NotTo(HaveOccurred())
		Expect(dnsServers).To(ContainElement("172.30.0.1"))

		mtu, err := winc.Property("winc_network/mtu")
		Expect(err).NotTo(HaveOccurred())
		Expect(mtu).To(Equal(1454))

		searchDomains, err := winc.Property("winc_network/search_domains")
		Expect(err).NotTo(HaveOccurred())
		Expect(searchDomains).To(Equal([]interface{}{"some-search-domain", "another-search-domain"}))
	})
})
