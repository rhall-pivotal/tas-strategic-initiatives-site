package manifest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("Apps Manager", func() {
	var instanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			instanceGroup = "control"
		} else {
			instanceGroup = "clock_global"
		}
	})

	Describe("BPM", func() {
		var appsManagerJobs []Job

		BeforeEach(func() {
			if productName == "srt" {
				appsManagerJobs = []Job{
					{
						InstanceGroup: "control",
						Name:          "push-apps-manager",
					},
				}
			} else {
				appsManagerJobs = []Job{
					{
						InstanceGroup: "clock_global",
						Name:          "push-apps-manager",
					},
				}
			}
		})

		It("co-locates the BPM job with all apps manager jobs", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, appsManagerJob := range appsManagerJobs {
				_, err = manifest.FindInstanceGroupJob(appsManagerJob.InstanceGroup, "bpm")
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("sets bpm.enabled to true", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			for _, appsManagerJob := range appsManagerJobs {
				manifestJob, err := manifest.FindInstanceGroupJob(appsManagerJob.InstanceGroup, appsManagerJob.Name)
				Expect(err).NotTo(HaveOccurred())

				Expect(manifestJob.Property("bpm/enabled")).To(BeTrue())
			}
		})
	})

	It("keeps the version number in docs link up-to-date", func() {
		manifest, err := product.RenderManifest(nil)
		Expect(err).NotTo(HaveOccurred())

		appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
		Expect(err).NotTo(HaveOccurred())

		Expect(appsManager.Property("apps_manager/white_labeling/nav_links")).To(ContainElement(HaveKeyWithValue(
			"href",
			MatchRegexp(`https://docs.pivotal.io/pivotalcf/\d+-\d+/pas/intro.html`),
		)))
	})

	Describe("Marketplace Url", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			appsManagerMarketplaceUrl, err := appsManager.Property("apps_manager/white_labeling/marketplace_url")
			Expect(err).NotTo(HaveOccurred())
			Expect(appsManagerMarketplaceUrl).To(BeNil())
		})

		Context("when the operator specifies a marketplace url", func() {
			It("applies it", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_marketplace_url": `custom-marketplace-url.com`,
				})
				Expect(err).NotTo(HaveOccurred())

				appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				appsManagerMarketplaceUrl, err := appsManager.Property("apps_manager/white_labeling/marketplace_url")
				Expect(err).NotTo(HaveOccurred())
				Expect(appsManagerMarketplaceUrl).To(Equal(`custom-marketplace-url.com`))
			})
		})
	})

	Describe("Secondary Navigation Links", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			navLinksDocsName, err := appsManager.Property("apps_manager/white_labeling/nav_links/0/name")
			Expect(err).NotTo(HaveOccurred())
			Expect(navLinksDocsName).To(Equal("Docs"))

			navLinksDocsHref, err := appsManager.Property("apps_manager/white_labeling/nav_links/0/href")
			Expect(err).NotTo(HaveOccurred())
			Expect(navLinksDocsHref).To(MatchRegexp(`https://docs.pivotal.io/pivotalcf/\d+-\d+/pas/intro.html`))

			navLinksToolsName, err := appsManager.Property("apps_manager/white_labeling/nav_links/1/name")
			Expect(err).NotTo(HaveOccurred())
			Expect(navLinksToolsName).To(Equal("Tools"))

			navLinksToolsHref, err := appsManager.Property("apps_manager/white_labeling/nav_links/1/href")
			Expect(err).NotTo(HaveOccurred())
			Expect(navLinksToolsHref).To(Equal("/tools"))
		})

		Context("when the operator specifies secondary navigation links", func() {
			It("applies them", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_nav_links": []map[string]interface{}{
						{
							"name": "custom-nav-1",
							"href": "custom-nav-1.com",
						},
						{
							"name": "custom-nav-2",
							"href": "custom-nav-2.com",
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())

				appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				navLinksCustomName1, err := appsManager.Property("apps_manager/white_labeling/nav_links/0/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(navLinksCustomName1).To(Equal(`custom-nav-1`))

				navLinksCustomHref1, err := appsManager.Property("apps_manager/white_labeling/nav_links/0/href")
				Expect(err).NotTo(HaveOccurred())
				Expect(navLinksCustomHref1).To(Equal("custom-nav-1.com"))

				navLinksCustomName2, err := appsManager.Property("apps_manager/white_labeling/nav_links/1/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(navLinksCustomName2).To(Equal(`custom-nav-2`))

				navLinksCustomHref2, err := appsManager.Property("apps_manager/white_labeling/nav_links/1/href")
				Expect(err).NotTo(HaveOccurred())
				Expect(navLinksCustomHref2).To(Equal("custom-nav-2.com"))
			})
		})
	})

	Describe("Memory", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			Expect(job.Property("apps_manager/memory")).To(BeNil())
			Expect(job.Property("invitations/memory")).To(BeNil())
		})

		Context("when the operator specifies memory limits", func() {
			It("applies them", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_memory":             1024,
					".properties.push_apps_manager_invitations_memory": 2048,
				})
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				Expect(job.Property("apps_manager/memory")).To(Equal(1024))
				Expect(job.Property("invitations/memory")).To(Equal(2048))
			})
		})
	})

	Describe("Polling intervals", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			Expect(job.Property("apps_manager/poll_interval")).To(Equal(30))
			Expect(job.Property("apps_manager/app_poll_interval")).To(Equal(10))
		})

		Context("when the operator specifies a poll interval", func() {
			It("applies them", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_app_poll_interval": 333,
					".properties.push_apps_manager_poll_interval":     666,
				})
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				Expect(job.Property("apps_manager/poll_interval")).To(Equal(666))
				Expect(job.Property("apps_manager/app_poll_interval")).To(Equal(333))
			})
		})
	})

	Describe("Foundations", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			foundations, err := appsManager.Property("apps_manager/foundations")
			Expect(err).NotTo(HaveOccurred())
			Expect(foundations).To(BeNil())
		})

		Context("when the operator specifies foundations", func() {
			It("applies them", func() {
				manifest, err := product.RenderManifest(map[string]interface{}{
					".properties.push_apps_manager_foundations": `{"foundation1": {"ccUrl": "api.foundation.my-env.com"}}`,
				})
				Expect(err).NotTo(HaveOccurred())

				appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				foundations, err := appsManager.Property("apps_manager/foundations")
				Expect(err).NotTo(HaveOccurred())
				Expect(foundations).To(Equal(`{"foundation1": {"ccUrl": "api.foundation.my-env.com"}}`))
			})
		})
	})

	Describe("Identity Providers", func() {
		It("fetches the SAML providers", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.uaa":                                "saml",
				".properties.uaa.saml.display_name":              "Some Display Name",
				".properties.uaa.saml.sso_name":                  "Okta",
				".properties.uaa.saml.require_signed_assertions": true,
			})
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			samlProviders, err := appsManager.Property("login/saml/providers")
			Expect(err).NotTo(HaveOccurred())
			Expect(samlProviders).To(HaveKey("Okta"))
		})
	})

	Describe("Backup and Restore", func() {
		Context("on the backup_restore instance group", func() {
			It("wires configurable fields for bbr-apps-manager", func() {
				expectations := []struct{
					uiReference  string
					uiValue      interface{}
					propertyPath string
					match        types.GomegaMatcher
				}{
					{".cloud_controller.system_domain", "example.com", "cf/api_url", Equal("https://api.example.com")},
					{".cloud_controller.system_domain", "example.com", "cf/uaa_url", Equal("https://login.example.com")},
					{".cloud_controller.system_domain", "example.com", "cf/notifications_service_url", Equal("https://notifications.example.com")},
					{".cloud_controller.system_domain", "example.com", "cf/system_domain", Equal("example.com")},
					{".cloud_controller.apps_domain", "example.com", "cf/apps_domain", Equal("example.com")},
					{".ha_proxy.skip_cert_verify", true, "ssl/skip_cert_verify", BeTrue()},
					{".properties.cf_dial_timeout_in_seconds", 10, "apps_manager/cf_dial_timeout", Equal(10)},
				}

				uiConfig := make(map[string]interface{})
				for _, expectation := range expectations {
					uiConfig[expectation.uiReference] = expectation.uiValue
				}

				manifest, err := product.RenderManifest(uiConfig)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				for _, expectation := range expectations {
					Expect(job.Property(expectation.propertyPath)).To(expectation.match)
				}
			})

			It("wires non-configurable fields for bbr-apps-manager", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "bbr-apps-manager")
				Expect(err).NotTo(HaveOccurred())

				_, err = job.Property("cf/admin_username")
				Expect(err).NotTo(HaveOccurred())

				_, err = job.Property("cf/admin_password")
				Expect(err).NotTo(HaveOccurred())
			})

			It("templates the push-apps-manager job", func() {
				manifest, err := product.RenderManifest(nil)
				Expect(err).NotTo(HaveOccurred())

				By("templating the push-apps-manager", func() {
					_, err := manifest.FindInstanceGroupJob("backup_restore", "push-apps-manager")
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})
	})

	Describe("Networking", func() {
		It("uses the spec defaults", func() {
			manifest, err := product.RenderManifest(nil)
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			networkingSelfService, err := appsManager.Property("networking/enable_space_developer_self_service")
			Expect(err).NotTo(HaveOccurred())
			Expect(networkingSelfService).To(BeFalse())
		})

		It("fetches the state of the networking space developer self-service checkbox", func() {
			manifest, err := product.RenderManifest(map[string]interface{}{
				".properties.cf_networking_enable_space_developer_self_service": true,
			})
			Expect(err).NotTo(HaveOccurred())

			appsManager, err := manifest.FindInstanceGroupJob(instanceGroup, "push-apps-manager")
			Expect(err).NotTo(HaveOccurred())

			networkingSelfService, err := appsManager.Property("networking/enable_space_developer_self_service")
			Expect(err).NotTo(HaveOccurred())
			Expect(networkingSelfService).To(BeTrue())
		})
	})
})
