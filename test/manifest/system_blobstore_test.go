package manifest_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("System Blobstore", func() {
	var ccInstanceGroup string
	var blobstoreInstanceGroup string

	BeforeEach(func() {
		if productName == "srt" {
			ccInstanceGroup = "control"
			blobstoreInstanceGroup = "blobstore"
		} else {
			ccInstanceGroup = "cloud_controller"
			blobstoreInstanceGroup = "nfs_server"
		}
	})

	Describe("internal blobstore", func() {
		It("configures the internal blobstore", func() {
			inputProperties := map[string]interface{}{
				".properties.system_blobstore":                                      "internal",
				".properties.system_blobstore_ccpackage_max_valid_packages_stored":  3,
				".properties.system_blobstore_ccdroplet_max_staged_droplets_stored": 3,
			}

			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			By("setting properties on cloud_controller_ng")
			job, err := manifest.FindInstanceGroupJob(ccInstanceGroup, "cloud_controller_ng")
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

			By("setting properties on blobstore")
			job, err = manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
			Expect(err).NotTo(HaveOccurred())

			selectDirectoriesToBackup, err := job.Property("select_directories_to_backup")
			Expect(err).NotTo(HaveOccurred())
			Expect(selectDirectoriesToBackup).To(ConsistOf("buildpacks", "packages", "droplets"))

			internalReleaseLevelBackup, err := job.Property("release_level_backup")
			Expect(err).NotTo(HaveOccurred())
			Expect(internalReleaseLevelBackup).To(BeTrue())

		})

		Context("when backup level skip_droplets is selected", func() {
			It("configures the select_directories_to_backup without droplets", func() {
				inputProperties := map[string]interface{}{
					".properties.system_blobstore_backup_level": "skip_droplets",
				}
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				By("setting properties on blobstore")
				job, err := manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				selectDirectoriesToBackup, err := job.Property("select_directories_to_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(selectDirectoriesToBackup).To(ConsistOf("buildpacks", "packages"))
			})
		})

		Context("when backup level skip_droplets_packages is selected", func() {
			It("configures the select_directories_to_backup without droplets", func() {
				inputProperties := map[string]interface{}{
					".properties.system_blobstore_backup_level": "skip_droplets_packages",
				}
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				By("setting properties on blobstore")
				job, err := manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				selectDirectoriesToBackup, err := job.Property("select_directories_to_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(selectDirectoriesToBackup).To(ConsistOf("buildpacks"))
			})
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

				By("disabling internal blobstore backup")
				job, err = manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				internalReleaseLevelBackup, err := job.Property("release_level_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalReleaseLevelBackup).To(BeFalse())

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

			When("backup level skip_droplets is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets"
				})

				It("only templates the buildpacks and packages buckets", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range []string{"buildpacks", "packages"} {
						_, err = job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
					}

					_, err = job.Property("buckets/droplets")
					Expect(err).To(HaveOccurred())
				})
			})

			When("backup level skip_droplets_packages is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets_packages"
				})

				It("only templates the buildpacks bucket", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					_, err = job.Property("buckets/buildpacks")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range []string{"packages", "droplets"} {
						_, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).To(HaveOccurred())
					}
				})
			})

			When("backup path style is enabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.path_style_s3_urls"] = "true"
				})

				It("templates force_path_style=true", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					value, err := job.Property("force_path_style")
					Expect(err).NotTo(HaveOccurred())

					Expect(value).To(BeTrue())
				})
			})

			When("backup path style is disabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.path_style_s3_urls"] = "false"
				})

				It("templates force_path_style=true", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-versioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					value, err := job.Property("force_path_style")
					Expect(err).NotTo(HaveOccurred())

					Expect(value).To(BeFalse())
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

				By("disabling internal blobstore backup")
				job, err = manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				internalReleaseLevelBackup, err := job.Property("release_level_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalReleaseLevelBackup).To(BeFalse())
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

			When("backup level skip_droplets is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets"
				})

				It("only templates the buildpacks and packages buckets", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range []string{"buildpacks", "packages"} {
						_, err = job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).NotTo(HaveOccurred())
					}

					_, err = job.Property("buckets/droplets")
					Expect(err).To(HaveOccurred())
				})
			})

			When("backup level skip_droplets_packages is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.iam_instance_profile_authentication"] = true
					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets_packages"
				})

				It("only templates the buildpacks bucket", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					_, err = job.Property("buckets/buildpacks")
					Expect(err).NotTo(HaveOccurred())

					for _, bucket := range []string{"packages", "droplets"} {
						_, err := job.Property(fmt.Sprintf("buckets/%s", bucket))
						Expect(err).To(HaveOccurred())
					}
				})
			})

			When("backup path style is enabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.path_style_s3_urls"] = "true"
				})

				It("templates force_path_style=true", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					value, err := job.Property("force_path_style")
					Expect(err).NotTo(HaveOccurred())

					Expect(value).To(BeTrue())
				})
			})

			When("backup path style is disabled", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external.path_style_s3_urls"] = "false"
				})

				It("templates force_path_style=true", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "s3-unversioned-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					value, err := job.Property("force_path_style")
					Expect(err).NotTo(HaveOccurred())

					Expect(value).To(BeFalse())
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

		Context("when soft delete is not configured", func() {
			It("disables the azure-blobstore-backup-restorer", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err := job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeFalse())

				By("disabling internal blobstore backup")
				job, err = manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				internalReleaseLevelBackup, err := job.Property("release_level_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalReleaseLevelBackup).To(BeFalse())
			})
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

			When("backup level skip_droplets is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external_azure.restore_from_account_name"] = "some-restore-account-name"
					inputProperties[".properties.system_blobstore.external_azure.restore_from_access_key"] = map[string]string{
						"secret": "some-restore-access-key",
					}

					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets"
				})

				It("only templates the buildpacks and packages containers", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					for _, container := range []string{"buildpacks", "packages"} {
						_, err = job.Property(fmt.Sprintf("containers/%s", container))
						Expect(err).NotTo(HaveOccurred())
					}

					_, err = job.Property("containers/droplets")
					Expect(err).To(HaveOccurred())
				})
			})

			When("backup level skip_droplets_packages is selected", func() {
				BeforeEach(func() {
					inputProperties[".properties.system_blobstore.external_azure.restore_from_account_name"] = "some-restore-account-name"
					inputProperties[".properties.system_blobstore.external_azure.restore_from_access_key"] = map[string]string{
						"secret": "some-restore-access-key",
					}

					inputProperties[".properties.system_blobstore_backup_level"] = "skip_droplets_packages"
				})

				It("only templates the buildpacks container", func() {
					manifest, err := product.RenderManifest(inputProperties)
					Expect(err).NotTo(HaveOccurred())

					job, err := manifest.FindInstanceGroupJob("backup_restore", "azure-blobstore-backup-restorer")
					Expect(err).NotTo(HaveOccurred())

					_, err = job.Property("containers/buildpacks")
					Expect(err).NotTo(HaveOccurred())

					for _, container := range []string{"packages", "droplets"} {
						_, err := job.Property(fmt.Sprintf("containers/%s", container))
						Expect(err).To(HaveOccurred())
					}
				})
			})
		})
	})

	Describe("gcs compatible with service account", func() {
		var inputProperties map[string]interface{}

		When("backup level all is selected", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_blobstore": "external_gcs_service_account",
					".properties.system_blobstore.external_gcs_service_account.buildpacks_bucket":        "some-buildpacks-bucket",
					".properties.system_blobstore.external_gcs_service_account.droplets_bucket":          "some-droplets-bucket",
					".properties.system_blobstore.external_gcs_service_account.packages_bucket":          "some-packages-bucket",
					".properties.system_blobstore.external_gcs_service_account.resources_bucket":         "some-resources-bucket",
					".properties.system_blobstore.external_gcs_service_account.service_account_json_key": "service-account-json-key",
					".properties.system_blobstore.external_gcs_service_account.project_id":               "dontcare",
					".properties.system_blobstore.external_gcs_service_account.service_account_email":    "dontcare",
					".properties.system_blobstore.external_gcs_service_account.backup_bucket":            "my-backup-bucket",
				}
			})

			It("enables the gcs-blobstore-backup-restorer", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "gcs-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err := job.Property("gcp_service_account_key")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(Equal("service-account-json-key"))

				job, err = manifest.FindInstanceGroupJob("backup_restore", "gcs-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				jobEnabled, err = job.Property("enabled")
				Expect(err).NotTo(HaveOccurred())
				Expect(jobEnabled).To(BeTrue())

				buildpacksBackupName, err := job.Property("buckets/buildpack/backup_bucket_name")
				Expect(err).NotTo(HaveOccurred())
				Expect(buildpacksBackupName).To(Equal("my-backup-bucket"))

				dropletsBackupName, err := job.Property("buckets/droplet/backup_bucket_name")
				Expect(err).NotTo(HaveOccurred())
				Expect(dropletsBackupName).To(Equal("my-backup-bucket"))

				packagesBackupName, err := job.Property("buckets/package/backup_bucket_name")
				Expect(err).NotTo(HaveOccurred())
				Expect(packagesBackupName).To(Equal("my-backup-bucket"))

				By("disabling internal blobstore backup")
				job, err = manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
				Expect(err).NotTo(HaveOccurred())

				internalReleaseLevelBackup, err := job.Property("release_level_backup")
				Expect(err).NotTo(HaveOccurred())
				Expect(internalReleaseLevelBackup).To(BeFalse())
			})

		})

		When("backup level skip_droplets is selected", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_blobstore": "external_gcs_service_account",
					".properties.system_blobstore.external_gcs_service_account.buildpacks_bucket":        "some-buildpacks-bucket",
					".properties.system_blobstore.external_gcs_service_account.droplets_bucket":          "some-droplets-bucket",
					".properties.system_blobstore.external_gcs_service_account.packages_bucket":          "some-packages-bucket",
					".properties.system_blobstore.external_gcs_service_account.resources_bucket":         "some-resources-bucket",
					".properties.system_blobstore.external_gcs_service_account.service_account_json_key": "service-account-json-key",
					".properties.system_blobstore.external_gcs_service_account.project_id":               "dontcare",
					".properties.system_blobstore.external_gcs_service_account.service_account_email":    "dontcare",
					".properties.system_blobstore.external_gcs_service_account.backup_bucket":            "my-backup-bucket",
					".properties.system_blobstore_backup_level":                                          "skip_droplets",
				}
			})

			It("only templates the buildpacks and packages buckets", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "gcs-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				for _, container := range []string{"buildpack", "package"} {
					_, err = job.Property(fmt.Sprintf("buckets/%s", container))
					Expect(err).NotTo(HaveOccurred())
				}

				_, err = job.Property("buckets/droplet")
				Expect(err).To(HaveOccurred())
			})
		})

		When("backup level skip_droplets_packages is selected", func() {
			BeforeEach(func() {
				inputProperties = map[string]interface{}{
					".properties.system_blobstore": "external_gcs_service_account",
					".properties.system_blobstore.external_gcs_service_account.buildpacks_bucket":        "some-buildpacks-bucket",
					".properties.system_blobstore.external_gcs_service_account.droplets_bucket":          "some-droplets-bucket",
					".properties.system_blobstore.external_gcs_service_account.packages_bucket":          "some-packages-bucket",
					".properties.system_blobstore.external_gcs_service_account.resources_bucket":         "some-resources-bucket",
					".properties.system_blobstore.external_gcs_service_account.service_account_json_key": "service-account-json-key",
					".properties.system_blobstore.external_gcs_service_account.project_id":               "dontcare",
					".properties.system_blobstore.external_gcs_service_account.service_account_email":    "dontcare",
					".properties.system_blobstore.external_gcs_service_account.backup_bucket":            "my-backup-bucket",
					".properties.system_blobstore_backup_level":                                          "skip_droplets_packages",
				}
			})

			It("only templates the buildpacks bucket", func() {
				manifest, err := product.RenderManifest(inputProperties)
				Expect(err).NotTo(HaveOccurred())

				job, err := manifest.FindInstanceGroupJob("backup_restore", "gcs-blobstore-backup-restorer")
				Expect(err).NotTo(HaveOccurred())

				_, err = job.Property("buckets/buildpack")
				Expect(err).NotTo(HaveOccurred())

				for _, container := range []string{"package", "droplet"} {
					_, err := job.Property(fmt.Sprintf("buckets/%s", container))
					Expect(err).To(HaveOccurred())
				}
			})
		})
	})

	Describe("gcs compatible with access key and secret key", func() {
		var inputProperties map[string]interface{}

		BeforeEach(func() {
			inputProperties = map[string]interface{}{
				".properties.system_blobstore":                                "external_gcs",
				".properties.system_blobstore.external_gcs.buildpacks_bucket": "foo",
				".properties.system_blobstore.external_gcs.droplets_bucket":   "foo",
				".properties.system_blobstore.external_gcs.packages_bucket":   "foo",
				".properties.system_blobstore.external_gcs.resources_bucket":  "foo",
				".properties.system_blobstore.external_gcs.access_key":        "foo",
				".properties.system_blobstore.external_gcs.secret_key": map[string]string{
					"secret": "some-access-key",
				},
			}
		})

		It("disables internal blobstore backup", func() {
			manifest, err := product.RenderManifest(inputProperties)
			Expect(err).NotTo(HaveOccurred())

			job, err := manifest.FindInstanceGroupJob(blobstoreInstanceGroup, "blobstore")
			Expect(err).NotTo(HaveOccurred())

			internalReleaseLevelBackup, err := job.Property("release_level_backup")
			Expect(err).NotTo(HaveOccurred())
			Expect(internalReleaseLevelBackup).To(BeFalse())
		})

	})
})
