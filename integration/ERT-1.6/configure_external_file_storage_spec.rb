require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.6.X External File Storage', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings.ops_manager.username,
      password: env_settings.ops_manager.password,
    )
  end

  it 'configures the external system blobstore' do
    file_storage_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('system_blobstore')
    file_storage_form.open_form
    file_storage_form.fill_in_selector_property(
      selector_input_reference: '.properties.system_blobstore',
      selector_name: 'external',
      selector_value: 'external',
      sub_field_answers: {
        '.properties.system_blobstore.external.endpoint' => {
          attribute_value: elastic_runtime_settings.file_storage.endpoint,
        },
        '.properties.system_blobstore.external.access_key' => {
          attribute_value: elastic_runtime_settings.file_storage.access_key,
        },
        '.properties.system_blobstore.external.secret_key' => {
          attribute_name: 'secret',
          attribute_value: elastic_runtime_settings.file_storage.secret_key,
        },
        '.properties.system_blobstore.external.buildpacks_bucket' => {
          attribute_value: elastic_runtime_settings.file_storage.buildpacks_bucket,
        },
        '.properties.system_blobstore.external.droplets_bucket' => {
          attribute_value: elastic_runtime_settings.file_storage.droplets_bucket,
        },
        '.properties.system_blobstore.external.packages_bucket' => {
          attribute_value: elastic_runtime_settings.file_storage.packages_bucket,
        },
        '.properties.system_blobstore.external.resources_bucket' => {
          attribute_value: elastic_runtime_settings.file_storage.resources_bucket,
        },
      },
    )
    file_storage_form.save_form
  end

  it 'scales nfs down to zero' do
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings.name)
    resource_config.set_instances_for_job('nfs_server', 0)
  end
end
