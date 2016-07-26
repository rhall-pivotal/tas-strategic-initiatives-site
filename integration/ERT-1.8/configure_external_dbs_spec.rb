require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.8.X External Databases', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'configures the external databases' do
    database_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('system_database')
    database_form.open_form
    database_form.fill_in_selector_property(
      selector_input_reference: '.properties.system_database',
      selector_name: 'external',
      selector_value: 'external',
      sub_field_answers: {
        '.properties.system_database.external.host' => {
          attribute_value: elastic_runtime_settings['rds']['host'],
        },
        '.properties.system_database.external.port' => {
          attribute_value: elastic_runtime_settings['rds']['port'],
        },
        '.properties.system_database.external.username' => {
          attribute_value: elastic_runtime_settings['rds']['username'],
        },
        '.properties.system_database.external.password' => {
          attribute_name: 'secret',
          attribute_value: elastic_runtime_settings['rds']['password'],
        },
      },
    )
    database_form.save_form
  end

  it 'scales down internal resources to zero' do
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_instances_for_job('uaadb', 0)
    resource_config.set_instances_for_job('ccdb', 0)
  end
end
