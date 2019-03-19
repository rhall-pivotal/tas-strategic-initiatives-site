package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("winc_network", func() {
	It("sets dns_servers and mtu", func() {
		manifest, err := product.RenderManifest(map[string]interface{}{})

		winc, err := manifest.FindInstanceGroupJob("windows_diego_cell", "winc-network-hns-acls")
		Expect(err).NotTo(HaveOccurred())

		dnsServers, err := winc.Property("winc_network/dns_servers")
		Expect(dnsServers).To(ContainElement("172.30.0.1"))

		mtu, err := winc.Property("winc_network/mtu")
		Expect(mtu).To(Equal(1454))
	})
})
