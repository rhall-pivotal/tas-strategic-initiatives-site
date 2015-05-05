require_relative '../opsman_ui_helper'

RSpec.describe 'Configuring µBosh', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:test_settings) { fetch_test_settings }

  it 'creates the admin user and logs in' do
    poll_up_to_mins(10) do
      current_ops_manager.setup_page.setup_or_login(
        user: test_settings.ops_manager.username,
        password: test_settings.ops_manager.password,
      )

      expect(page).to have_content('Installation Dashboard')
    end
  end

  it 'configures µBosh' do
    current_ops_manager.ops_manager_director.configure_microbosh(test_settings)
  end
end
