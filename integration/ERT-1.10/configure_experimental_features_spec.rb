require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.10.X Experimental Features', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'enables all of the experimental features' do
    experimental_features_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('experimental_features')
    experimental_features_form.open_form

    all('input[type=checkbox]').each do |checkbox|
      checkbox.click unless checkbox.checked?
    end

    experimental_features_form.save_form
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])

    case env_settings['iaas_type']
    when 'aws'
      resource_config.set_elb_names_for_job('tcp_router', elastic_runtime_settings['tcp_elb_name'])
    when 'openstack'
      resource_config.set_floating_ips_for_job('tcp_router', elastic_runtime_settings['tcp_router_floating_ips'])
    end
  end
end
