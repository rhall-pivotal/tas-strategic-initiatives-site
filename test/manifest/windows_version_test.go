package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal-cf/planitest"
)

var _ = Describe("WindowsVersion", func() {
	var (
		manifest planitest.Manifest
		err      error
	)

	BeforeEach(func() {
		manifest, err = product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("uses 1803 as its Windows Version", func() {
		By("having a winc-network-1803 job on the diego cell", func() {
			_, err = manifest.FindInstanceGroupJob("windows_diego_cell", "winc-network-1803")
			Expect(err).NotTo(HaveOccurred())
		})

		By("using the path to winc-network-1803.exe for the garden-windows network_plugin", func() {
			job, err := manifest.FindInstanceGroupJob("windows_diego_cell", "garden-windows")
			Expect(err).NotTo(HaveOccurred())

			plugin, err := job.Property("garden/network_plugin")
			Expect(err).NotTo(HaveOccurred())
			Expect(plugin).To(Equal("/var/vcap/packages/winc-network-1803/winc-network.exe"))
		})

		By("configures the network_plugin_extra_args for garden-windows", func() {
			job, err := manifest.FindInstanceGroupJob("windows_diego_cell", "garden-windows")
			Expect(err).NotTo(HaveOccurred())

			args, err := job.Property("garden/network_plugin_extra_args")
			Expect(err).NotTo(HaveOccurred())
			Expect(args).To(ContainElement("--configFile=/var/vcap/jobs/winc-network-1803/config/interface.json"))
			Expect(args).To(ContainElement("--log=/var/vcap/sys/log/winc-network-1803/winc-network.log"))
		})

		By("configuring the windows2016 preloaded_rootfs to point to /var/vcap/packages/windows1803fs", func() {
			job, err := manifest.FindInstanceGroupJob("windows_diego_cell", "rep_windows")
			Expect(err).NotTo(HaveOccurred())

			args, err := job.Property("diego/rep/preloaded_rootfses")
			Expect(err).NotTo(HaveOccurred())
			Expect(args).To(ContainElement("windows2016:oci:///C:/var/vcap/packages/windows1803fs"))
		})

		By("configuring the groot cached_image_uris to point to /var/vcap/packages/windows1803fs", func() {
			job, err := manifest.FindInstanceGroupJob("windows_diego_cell", "groot")
			Expect(err).NotTo(HaveOccurred())

			args, err := job.Property("groot/cached_image_uris")
			Expect(err).NotTo(HaveOccurred())
			Expect(args).To(ContainElement("oci:///C:/var/vcap/packages/windows1803fs"))
		})

		By("having a windows1803fs job on the diego cell", func() {
			_, err = manifest.FindInstanceGroupJob("windows_diego_cell", "windows1803fs")
			Expect(err).NotTo(HaveOccurred())
		})

	})

})
