require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.6.X Experimental Features', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings.ops_manager.username,
      password: env_settings.ops_manager.password,
    )
  end

  it 'enables all of the experimental features' do
    experimental_features_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('experimental_features')
    experimental_features_form.open_form

    all('input[type=checkbox]').each do |checkbox|
      checkbox.click unless checkbox.checked?
    end

    experimental_features_form.save_form
  end

  it 'enables the diego features' do
    diego_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('diego')

    diego_form.open_form

    diego_form.property('.cloud_controller.default_to_diego_backend').set(true)

    diego_form.save_form
  end

  it 'disables the errands that do not work with diego' do
    %w().each do |errand|
      current_ops_manager.product(elastic_runtime_settings.name)
        .product_errands.disable_errand(errand)
    end
  end
end
