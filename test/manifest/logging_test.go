package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging", func() {
	Describe("traffic controller", func() {
		It("disables support for forwarding syslog to metron", func() {
			err := product.Configure(map[string]interface{}{})
			Expect(err).NotTo(HaveOccurred())

			manifest, err := product.RenderManifest()
			Expect(err).NotTo(HaveOccurred())

			agent, err := manifest.FindInstanceGroupJob("loggregator_trafficcontroller", "metron_agent")
			Expect(err).NotTo(HaveOccurred())

			syslogForwardingEnabled, err := agent.Property("syslog_daemon_config/enable")
			Expect(syslogForwardingEnabled).To(BeFalse())
		})
	})
})
