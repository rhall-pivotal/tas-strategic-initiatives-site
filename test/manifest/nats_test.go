package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NATS", func() {
	Describe("Container networking", func() {
		var (
			inputProperties   map[string]interface{}
			natsInstanceGroup string
		)

		BeforeEach(func() {
			if productName == "ert" {
				natsInstanceGroup = "nats"
			} else {
				natsInstanceGroup = "database"
			}
			inputProperties = map[string]interface{}{}
		})

		Describe("NATS", func() {
			It("enabled internal tls", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats")
				Expect(err).NotTo(HaveOccurred())

				username, err := job.Property("nats/internal/tls/enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(username).To(BeTrue())
			})

			It("has credentials", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats")
				Expect(err).NotTo(HaveOccurred())

				username, err := job.Property("nats/user")
				Expect(err).NotTo(HaveOccurred())
				Expect(username).To(Equal("((nats-credentials.username))"))

				password, err := job.Property("nats/password")
				Expect(err).NotTo(HaveOccurred())
				Expect(password).To(Equal("((nats-credentials.password))"))
			})

			It("has certs for internal nats cluster connections", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats")
				Expect(err).NotTo(HaveOccurred())

				ca, err := job.Property("nats/internal/tls/ca")
				Expect(err).NotTo(HaveOccurred())
				Expect(ca).To(Equal("fake-ops-manager-ca-certificate"))

				internalCert, err := job.Property("nats/internal/tls/certificate")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalCert).To(BeNil())

				internalKey, err := job.Property("nats/internal/tls/private_key")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalKey).To(BeNil())
			})

			It("has a hostname", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats")
				Expect(err).NotTo(HaveOccurred())

				hostname, err := job.Property("nats/hostname")
				Expect(err).NotTo(HaveOccurred())
				Expect(hostname).To(Equal("nats.service.cf.internal"))
			})
		})

		Describe("NATS TLS", func() {
			It("enabled internal tls", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats-tls")
				Expect(err).NotTo(HaveOccurred())

				username, err := job.Property("nats/internal/tls/enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(username).To(BeTrue())
			})

			It("has credentials", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats-tls")
				Expect(err).NotTo(HaveOccurred())

				username, err := job.Property("nats/user")
				Expect(err).NotTo(HaveOccurred())
				Expect(username).To(Equal("((nats-credentials.username))"))

				password, err := job.Property("nats/password")
				Expect(err).NotTo(HaveOccurred())
				Expect(password).To(Equal("((nats-credentials.password))"))
			})

			It("has certs for external clients", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats-tls")
				Expect(err).NotTo(HaveOccurred())

				ca, err := job.Property("nats/external/tls/ca")
				Expect(err).NotTo(HaveOccurred())
				Expect(ca).To(Equal("fake-ops-manager-ca-certificate"))

				externalCert, err := job.Property("nats/external/tls/certificate")
				Expect(err).NotTo(HaveOccurred())
				Expect(externalCert).To(BeNil())

				externalKey, err := job.Property("nats/external/tls/private_key")
				Expect(err).NotTo(HaveOccurred())
				Expect(externalKey).To(BeNil())
			})

			It("has certs for internal nats cluster connections", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats-tls")
				Expect(err).NotTo(HaveOccurred())

				ca, err := job.Property("nats/internal/tls/ca")
				Expect(err).NotTo(HaveOccurred())
				Expect(ca).To(Equal("fake-ops-manager-ca-certificate"))

				internalCert, err := job.Property("nats/internal/tls/certificate")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalCert).To(BeNil())

				internalKey, err := job.Property("nats/internal/tls/private_key")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalKey).To(BeNil())
			})

			It("has a hostname", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(natsInstanceGroup, "nats-tls")
				Expect(err).NotTo(HaveOccurred())

				hostname, err := job.Property("nats/hostname")
				Expect(err).NotTo(HaveOccurred())
				Expect(hostname).To(Equal("nats.service.cf.internal"))
			})
		})
	})
})
