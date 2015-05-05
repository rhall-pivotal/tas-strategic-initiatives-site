require_relative '../opsman_ui_helper'

RSpec.describe 'Install Elastic Runtime 1.5.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:test_settings) { fetch_test_settings }

  let(:elastic_runtime_settings) { test_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: test_settings.ops_manager.username,
      password: test_settings.ops_manager.password,
    )
  end

  it 'installs runtime' do
    current_ops_manager.product_dashboard.apply_updates
  end

  it 'waits for a successful installation' do
    poll_up_to_mins(360) do
      expect(current_ops_manager.state_change_progress).to be_state_change_success
    end
  end
end
