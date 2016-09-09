require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.5.X HA Instance Counts', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'sets the instance counts' do
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_instances_for_job('nats', 2)
    resource_config.set_instances_for_job('etcd_server', 1)
    resource_config.set_instances_for_job('router', 2)
    resource_config.set_instances_for_job('mysql_proxy', 2)
    resource_config.set_instances_for_job('mysql', 3)
    resource_config.set_instances_for_job('cloud_controller', 2)
    resource_config.set_instances_for_job('ha_proxy', 2) unless env_settings['iaas_type'] == 'aws'
    resource_config.set_instances_for_job('health_manager', 2)
    resource_config.set_instances_for_job('cloud_controller_worker', 2)
    resource_config.set_instances_for_job('uaa', 2)
    resource_config.set_instances_for_job('dea', 4)
    resource_config.set_instances_for_job('doppler', 2)
    resource_config.set_instances_for_job('loggregator_trafficcontroller', 2)
  end
end
