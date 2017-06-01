require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.8.X AWS Diego Cell Instance', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'sets chooses a larger machine for the diego cells' do
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_vm_type_for_job('diego_cell', 'r3.2xlarge') if env_settings['iaas_type'] == 'aws'
  end
end
