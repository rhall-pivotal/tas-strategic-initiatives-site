package manifest_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "uaa"
		}
	})

	Describe("database connection", func() {
		It("configures TLS to the internal database", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			tlsEnabled, err := job.Property("uaadb/tls_enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(tlsEnabled).To(BeTrue())
		})

		It("trusts the certificate provided by the server", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			caCerts, err := job.Property("uaa/ca_certs")
			Expect(err).NotTo(HaveOccurred())
			Expect(caCerts).NotTo(BeEmpty())
		})

		It("requires TLS 1.2", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			tlsProtocols, err := job.Property("uaadb/tls_protocols")
			Expect(err).NotTo(HaveOccurred())
			Expect(tlsProtocols).To(Equal("TLSv1.2"))
		})
	})

	Describe("route registration", func() {
		It("tags the emitted metrics", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			routeRegistrar, err := manifest.FindInstanceGroupJob(instanceGroup, "route_registrar")
			Expect(err).NotTo(HaveOccurred())

			routes, err := routeRegistrar.Property("route_registrar/routes")
			Expect(err).ToNot(HaveOccurred())
			Expect(routes).To(ContainElement(HaveKeyWithValue("tags", map[interface{}]interface{}{
				"component": "uaa",
			})))
		})
	})

	Describe("BPM", func() {
		It("co-locates the BPM job with all diego jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			_, err = manifest.FindInstanceGroupJob(instanceGroup, "bpm")
			Expect(err).NotTo(HaveOccurred())
		})

		It("sets bpm.enabled to true", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			manifestJob, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			bpmEnabled, err := manifestJob.Property("bpm/enabled")
			Expect(err).NotTo(HaveOccurred())

			Expect(bpmEnabled).To(BeTrue())
		})
	})

	Describe("Clients", func() {
		It("apps_metrics has the expected permission scopes", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsScopes, err := uaa.Property("uaa/clients/apps_metrics/scope")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsScopes).To(Equal("cloud_controller.admin,cloud_controller.read,metrics.read,cloud_controller.admin_read_only"))

		})

		It("apps_metrics has the expected redirect uri", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics/redirect-uri")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com,https://metrics.sys.example.com/,https://metrics.sys.example.com/*,https://metrics-previous.sys.example.com,https://metrics-previous.sys.example.com/,https://metrics-previous.sys.example.com/*"))

		})

		It("apps_metrics_processing has the expected permission scopes", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsProcessingScopes, err := uaa.Property("uaa/clients/apps_metrics_processing/scope")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsProcessingScopes).To(Equal("openid,oauth.approvals,doppler.firehose,cloud_controller.admin,cloud_controller.admin_read_only"))

		})

		It("apps_metrics_processing has the expected redirect uri", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			appMetricsRedirectUri, err := uaa.Property("uaa/clients/apps_metrics_processing/redirect-uri")
			Expect(err).ToNot(HaveOccurred())
			Expect(appMetricsRedirectUri).To(Equal("https://metrics.sys.example.com,https://metrics-previous.sys.example.com"))

		})

		It("apps_manager_js client includes network.write and network.admin", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			rawScopes, err := uaa.Property("uaa/clients/apps_manager_js/scope")
			Expect(err).ToNot(HaveOccurred())

			scopes := strings.Split(rawScopes.(string), ",")
			Expect(scopes).To(ContainElement("network.write"))
			Expect(scopes).To(ContainElement("network.admin"))

			autoapproveList, err := uaa.Property("uaa/clients/apps_manager_js/autoapprove")
			Expect(err).ToNot(HaveOccurred())

			Expect(autoapproveList).To(ContainElement("network.write"))
			Expect(autoapproveList).To(ContainElement("network.admin"))
		})

		It("credhub_admin_client has credhub.read and credhub.write", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			uaa, err := manifest.FindInstanceGroupJob(instanceGroup, "uaa")
			Expect(err).NotTo(HaveOccurred())

			id, err := uaa.Property("uaa/clients/credhub_admin_client/id")
			Expect(err).ToNot(HaveOccurred())
			Expect(id).To(Equal("credhub_admin_client"))

			rawAuthorities, err := uaa.Property("uaa/clients/credhub_admin_client/authorities")
			Expect(err).ToNot(HaveOccurred())

			authorities := strings.Split(rawAuthorities.(string), ",")
			Expect(authorities).To(ConsistOf([]string{"credhub.read", "credhub.write"}))

			authorizedGrantTypes, err := uaa.Property("uaa/clients/credhub_admin_client/authorized-grant-types")
			Expect(err).ToNot(HaveOccurred())
			Expect(authorizedGrantTypes).To(Equal("client_credentials"))
		})

	})
})
