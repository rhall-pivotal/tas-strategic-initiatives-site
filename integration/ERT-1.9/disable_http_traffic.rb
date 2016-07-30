require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Disable HTTP Traffic in Elastic Runtime 1.9.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  it 'disables HTTP traffic to the HAProxy and UAA' do
    networking_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('networking')
    networking_form.open_form

    check 'networking[.properties.networking_point_of_entry][haproxy][.properties.networking_point_of_entry.haproxy.disable_http]'

    networking_form.save_form
  end
end
