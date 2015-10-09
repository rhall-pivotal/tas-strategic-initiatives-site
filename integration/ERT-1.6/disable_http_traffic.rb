require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Disable HTTP Traffic in Elastic Runtime 1.6.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings.ops_manager.username,
      password: env_settings.ops_manager.password,
    )
  end

  it 'disables HTTP traffic to the HAProxy and UAA' do
    security_config_form =
        current_ops_manager.product(elastic_runtime_settings.name).product_form('security_config')
    security_config_form.open_form

    check 'security_config[.ha_proxy.disable_http]'
    check 'security_config[.uaa.disable_http]'

    security_config_form.save_form
  end
end
