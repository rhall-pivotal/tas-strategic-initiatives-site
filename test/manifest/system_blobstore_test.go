package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("System Blobstore", func() {
	Describe("internal blobstore", func() {
		var instanceGroup string

		BeforeEach(func() {
			if productName == "srt" {
				instanceGroup = "control"
			} else {
				instanceGroup = "cloud_controller"
			}
		})

		It("configures the internal blobstore", func() {
			inputProperties := map[string]interface{}{
				".properties.system_blobstore_ccpackage_max_valid_packages_stored":  3,
				".properties.system_blobstore_ccdroplet_max_staged_droplets_stored": 3,
			}

			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			By("setting properties on cloud_controller_ng")
			job, err := manifest.FindInstanceGroupJob(instanceGroup, "cloud_controller_ng")
			Expect(err).NotTo(HaveOccurred())

			maxValidPackages, err := job.Property("cc/packages/max_valid_packages_stored")
			Expect(err).NotTo(HaveOccurred())
			Expect(maxValidPackages).To(Equal(3))

			maxStagedDroplets, err := job.Property("cc/droplets/max_staged_droplets_stored")
			Expect(err).NotTo(HaveOccurred())
			Expect(maxStagedDroplets).To(Equal(3))

			By("not enabling unversioned S3 backups", func() {})
			bbr, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
			Expect(err).NotTo(HaveOccurred())

			jobEnabled, err := bbr.Property("enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(jobEnabled).To(BeFalse())
		})
	})

	Describe("s3 compatible", func() {
		var (
			inputProperties map[string]interface{}
			backupBuckets   = []string{"buildpacks", "droplets", "packages"}
		)

		BeforeEach(func() {
			inputProperties = map[string]interface{}{
				".properties.system_blobstore":                            "external",
				".properties.system_blobstore.external.buildpacks_bucket": "some-buildpacks-bucket",
				".properties.system_blobstore.external.droplets_bucket":   "some-droplets-bucket",
				".properties.system_blobstore.external.packages_bucket":   "some-packages-bucket",
				".properties.system_blobstore.external.resources_bucket":  "some-resources-bucket",
			}
		})

		Context("when the user enables versioning", func() {
			It("enables the s3-versioned-blobstore-backup-restorer, and disables the s3-unversioned-blobstore-backup-restorer", func() {
				inputProperties[".properties.system_blobstore.external.versioning"] = true

				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err := job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeTrue())

				job, err = manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err = job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeFalse())
			})

			Context("and IAM instance profiles are disabled", func() {

				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = false
					inputProperties[".properties.system_blobstore.external.access_key"] = "some-access-key-id"
					inputProperties[".properties.system_blobstore.external.secret_key"] = map[string]string{
						"secret": "some-secret-access-key",
					}
				})

				It("specifies that backups use the provided access key", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range backupBuckets {
						bucketProperties, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(bucketProperties).NotTo(HaveKey("use_iam_profile"))
						Expect(bucketProperties).To(HaveKeyWithValue("aws_access_key_id", "some-access-key-id"))
						Expect(bucketProperties).To(HaveKeyWithValue("aws_secret_access_key", ContainSubstring("system_blobstore/external/secret_key")))
					}
				})
			})

			Context("and IAM instance profiles are enabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
				})

				It("specifies that backups should use it", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range backupBuckets {
						iamInstanceProfileAuthentication, err := job.Property(fmt.Sprintf("buckets/%s/use_iam_profile", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(iamInstanceProfileAuthentication).To(BeTrue())

						bucketProperties, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(bucketProperties).NotTo(HaveKey("aws_access_key_id"))
						Expect(bucketProperties).NotTo(HaveKey("aws_access_secret_key"))
					}
				})
			})
		})

		Context("when the user disables versioning", func() {
			BeforeEach(func() {
				inputProperties[".properties.system_blobstore.external.backup_region"] = "some-backup-region"
				inputProperties[".properties.system_blobstore.external.buildpacks_backup_bucket"] = "some-buildpacks-bucket"
				inputProperties[".properties.system_blobstore.external.droplets_backup_bucket"] = "some-droplets-bucket"
				inputProperties[".properties.system_blobstore.external.packages_backup_bucket"] = "some-packages-bucket"
			})

			It("disables the s3-versioned-blobstore-backup-restorer and enables the s3-unversioned-blobstore-backup-restorer", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err := job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeFalse())

				job, err = manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err = job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeTrue())

				buildpacksBackupRegion, err := job.Property("buckets/buildpacks/backup/region")
				Expect(err).NotTo(HaveOccurred())
				Expect(buildpacksBackupRegion).To(Equal("some-backup-region"))

				buildpacksBackupName, err := job.Property("buckets/buildpacks/backup/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(buildpacksBackupName).To(Equal("some-buildpacks-bucket"))

				dropletsBackupRegion, err := job.Property("buckets/droplets/backup/region")
				Expect(err).NotTo(HaveOccurred())
				Expect(dropletsBackupRegion).To(Equal("some-backup-region"))

				dropletsBackupName, err := job.Property("buckets/droplets/backup/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(dropletsBackupName).To(Equal("some-droplets-bucket"))

				packagesBackupRegion, err := job.Property("buckets/packages/backup/region")
				Expect(err).NotTo(HaveOccurred())
				Expect(packagesBackupRegion).To(Equal("some-backup-region"))

				packagesBackupName, err := job.Property("buckets/packages/backup/name")
				Expect(err).NotTo(HaveOccurred())
				Expect(packagesBackupName).To(Equal("some-packages-bucket"))
			})

			Context("and IAM instance profiles are disabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = false
					inputProperties[".properties.system_blobstore.external.access_key"] = "some-access-key-id"
					inputProperties[".properties.system_blobstore.external.secret_key"] = map[string]string{
						"secret": "some-secret-access-key",
					}
				})

				It("specifies that backups use the provided access key", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range backupBuckets {
						bucketProperties, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(bucketProperties).NotTo(HaveKey("use_iam_profile"))
						Expect(bucketProperties).To(HaveKeyWithValue("aws_access_key_id", "some-access-key-id"))
						Expect(bucketProperties).To(HaveKeyWithValue("aws_secret_access_key", ContainSubstring("system_blobstore/external/secret_key")))
					}
				})
			})

			Context("and IAM instance profiles are enabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
				})

				It("specifies that backups should use it", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range backupBuckets {
						iamInstanceProfileAuthentication, err := job.Property(fmt.Sprintf("buckets/%s/use_iam_profile", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(iamInstanceProfileAuthentication).To(BeTrue())

						bucketProperties, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
						Expect(bucketProperties).NotTo(HaveKey("aws_access_key_id"))
						Expect(bucketProperties).NotTo(HaveKey("aws_access_secret_key"))
					}
				})
			})
		})
	})

	Describe("azure storage", func() {
		var (
			inputProperties map[string]interface{}
			containers      = []string{"buildpacks", "droplets", "packages"}
		)

		BeforeEach(func() {
			inputProperties = map[string]interface{}{
				".properties.system_blobstore":                                     "external_azure",
				".properties.system_blobstore.external_azure.buildpacks_container": "some-buildpacks-bucket",
				".properties.system_blobstore.external_azure.droplets_container":   "some-droplets-bucket",
				".properties.system_blobstore.external_azure.packages_container":   "some-packages-bucket",
				".properties.system_blobstore.external_azure.resources_container":  "some-resources-bucket",
				".properties.system_blobstore.external_azure.account_name":         "some-account-name",
				".properties.system_blobstore.external_azure.access_key": map[string]string{
					"secret": "some-access-key",
				},
			}
		})

		It("disables the azure-blobstore-backup-restorer", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
			Expect(err).NotTo(HaveOccurred())

			jobEnabled, err := job.Property("enabled")
			Expect(err).NotTo(HaveOccurred())
			Expect(jobEnabled).To(BeFalse())
		})

		Context("when the user enables backup and restore", func() {
			BeforeEach(func() {
				inputProperties[".properties.system_blobstore.external_azure.enable_bbr"] = true
			})

			It("enables the azure-blobstore-backup-restorer", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err := job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeTrue())

				for _, container := range containers {
					containerProperties, err := job.Property(fmt.Sprintf("containers/%s", container))
					Expect(err).NotTo(HaveOccurred())
					Expect(containerProperties).To(HaveKeyWithValue("azure_storage_account", "some-account-name"))
					Expect(containerProperties).To(HaveKeyWithValue("azure_storage_key", ContainSubstring("system_blobstore/external_azure/access_key")))
				}
			})

			Context("with restore from credentials", func() {
				It("enables the azure-blobstore-backup-restorer", func() {
					inputProperties[".properties.system_blobstore.external_azure.restore_from_account_name"] = "some-restore-account-name"
					inputProperties[".properties.system_blobstore.external_azure.restore_from_access_key"] = map[string]string{
						"secret": "some-restore-access-key",
					}

					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					jobEnabled, err := job.Property("enabled")
					Expect(err).NotTo(HaveOccurred())
					Expect(jobEnabled).To(BeTrue())

					for _, container := range containers {
						restoreFromProperties, err := job.Property(fmt.Sprintf("containers/%s/restore_from", container))
						Expect(err).NotTo(HaveOccurred())
						Expect(restoreFromProperties).To(HaveKeyWithValue("azure_storage_account", "some-restore-account-name"))
						Expect(restoreFromProperties).To(HaveKeyWithValue("azure_storage_key", ContainSubstring("system_blobstore/external_azure/restore_from_access_key")))
					}
				})
			})
		})
	})
})
