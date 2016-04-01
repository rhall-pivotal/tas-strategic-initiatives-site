require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.6.X to use Postgres Databases', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'configures Database section to use Postgres' do
    database_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('system_database')
    database_form.open_form
    database_form.fill_in_selector_property(
      selector_input_reference: '.properties.system_database',
      selector_name: 'internal',
      selector_value: 'internal',
    )
    database_form.save_form
  end
end
