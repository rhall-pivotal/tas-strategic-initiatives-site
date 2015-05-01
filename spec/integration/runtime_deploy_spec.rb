require 'ops_manager_ui_drivers'
require 'capybara/webkit'
require 'yaml'
require 'recursive-open-struct'

RSpec.configure do |config|
  config.fail_fast = true

  Capybara.save_and_open_page_path = __dir__

  Capybara.configure do |c|
    c.default_driver = :webkit
  end

  config.after(:each) do |example|
    if example.exception
      page = save_page
      screenshot = save_screenshot(nil)

      exception = example.exception
      exception.define_singleton_method :message do
        super() +
          "\nHTML page: #{page}" +
          "\nScreenshot: #{screenshot}"
      end
    end
  end
end

describe 'Elastic Runtime and MySQL integration', order: :defined, type: :integration do
  include Capybara::DSL
  include OpsManagerUiDrivers::PageHelpers
  include OpsManagerUiDrivers::WaitHelper

  before(:all) do
    # Fixes following problem (since we use self-signed certs):
    #   Unable to load URL: https://localhost:5443/
    #   because of error loading https://localhost:5443/: Unknown error
    page.driver.browser.ignore_ssl_errors
    page.current_window.resize_to(1024, 1600) # avoid overlapping footer spec failures
  end

  let(:settings_hash) do
    YAML.load(<<YAML)
iaas_type: vsphere
ops_manager:
  url: https://pcf.ridge.cf-app.com
  username: pivotalcf
  password: pivotalcf
  ntp_servers: time1.sf.cf-app.com
  vcenter:
    creds:
      ip: 10.84.32.61
      username: root
      password: look74Beer
    datacenter: private
    datastore: public
    microbosh_vm_folder: ridge_vms
    microbosh_template_folder: ridge_templates
    microbosh_disk_path: ridge_disk
  availability_zones:
  - name: first-az
    cluster: geography
    resource_pool: ridge
  networks:
  - name: first-network
    identifier: ridge
    subnet: 10.85.88.0/24
    gateway: 10.85.88.1
    reserved_ips: 10.85.88.0-10.85.88.90
    dns: 10.87.8.10,10.87.8.11
  elastic_runtime:
    name: cf
    domain: ridge.cf-app.com
  ha_proxy_static_ips: 10.85.88.254
YAML
  end
  let(:test_settings) { RecursiveOpenStruct.new(settings_hash) }
  let(:current_ops_manager) { om_1_4(test_settings.ops_manager.url) }

  describe 'initial setup' do
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

  context 'with Elastic Runtime' do
    let(:elastic_runtime_settings) { test_settings.ops_manager.elastic_runtime }
    let(:product_file_path) { '/Users/pivotal/workspace/p-runtime/cf-1.5.0.0.alpha.726.bccd149.dirty.pivotal'}

    it 'uploads the product' do
      # current_ops_manager.product_dashboard.import_product_from(product_file_path)
    end

    it 'adds the product' do
      current_ops_manager.available_products.add_product_to_install(elastic_runtime_settings.name)
    end

    it 'configures the product' do
      ips_and_ports_form = current_ops_manager.product(elastic_runtime_settings.name).product_form('ha_proxy')
      ips_and_ports_form.open_form
      ips_and_ports_form.property('.ha_proxy.static_ips').set(elastic_runtime_settings.ha_proxy_static_ips)
      ips_and_ports_form.save_form

      security_config_form = current_ops_manager.product(elastic_runtime_settings.name).product_form('security_config')
      security_config_form.open_form
      security_config_form.property('.ha_proxy.skip_cert_verify').set(true)
      security_config_form.generate_self_signed_cert("*.#{elastic_runtime_settings.domain}")
      security_config_form.save_form

      cloud_controller_form =
        current_ops_manager.product(elastic_runtime_settings.name).product_form('cloud_controller')
      cloud_controller_form.open_form
      cloud_controller_form.property('.cloud_controller.system_domain').set(elastic_runtime_settings.domain)
      cloud_controller_form.property('.cloud_controller.apps_domain').set(elastic_runtime_settings.domain)
      cloud_controller_form.save_form
    end
  end

  describe 'applying changes' do
    it 'works once' do
      current_ops_manager.product_dashboard.apply_updates

      poll_up_to_mins(360) do
        expect(current_ops_manager.state_change_progress).to be_state_change_success
      end
    end
  end

  describe 'deleting the entire installation' do
    it 'succeeds' do
      current_ops_manager.product_dashboard.delete_whole_installation

      expect(current_ops_manager.product_dashboard).not_to be_delete_installation_available

      poll_up_to_mins(120) do
        expect(current_ops_manager.state_change_progress).to be_state_change_success
      end
    end
  end
end
