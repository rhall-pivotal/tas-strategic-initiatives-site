require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Upload/Add Elastic Runtime 1.4.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:test_settings) { fetch_test_settings }

  let(:elastic_runtime_settings) { test_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: test_settings.ops_manager.username,
      password: test_settings.ops_manager.password,
    )
  end

  it 'uploads the product' do
    current_ops_manager.product_dashboard.import_product_from(product_path)
  end

  it 'adds the product' do
    current_ops_manager.available_products.add_product_to_install(elastic_runtime_settings.name)
  end
end
