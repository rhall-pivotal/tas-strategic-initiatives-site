PROJECT_ROOT = File.expand_path(File.join(__dir__, '..'))

require 'rspec'
require 'ops_manager_ui_drivers'
require 'recursive-open-struct'
require 'capybara/webkit'
require 'yaml'

module SettingsHelper
  def environments_dir
    ENV.fetch('ENVIRONMENTS_DIR', File.join(PROJECT_ROOT, 'config', 'environments'))
  end

  def product_path
    ENV.fetch('PRODUCT_PATH')
  end

  def om_version
    ENV.fetch('OM_VERSION')
  end

  def ops_manager_driver
    case om_version
    when '1.4'
      om_1_4(fetch_test_settings.ops_manager.url)
    when '1.5'
      om_rc(fetch_test_settings.ops_manager.url)
    else
      fail "Unsupported Ops Manager Version #{om_version.inspect}"
    end
  end

  def environment_name
    ENV.fetch('ENVIRONMENT_NAME')
  end

  def fetch_test_settings
    RecursiveOpenStruct.new(
      YAML.load_file(File.join(environments_dir, "#{environment_name}.yml"))
    )
  end
end

Capybara.save_and_open_page_path = File.join(__dir__, 'tmp')

Capybara.configure do |c|
  c.default_driver = :webkit
end

RSpec.configure do |config|
  config.fail_fast = true
  config.default_formatter = 'documentation'

  config.include(SettingsHelper)
  config.include(Capybara::DSL)
  config.include(OpsManagerUiDrivers::PageHelpers)
  config.include(OpsManagerUiDrivers::WaitHelper)

  config.before(:all) do
    page.driver.browser.ignore_ssl_errors # allows use of self-signed SSL certificates
    page.current_window.resize_to(1024, 1600) # avoid overlapping footer spec failures
  end

  config.after(:each) do |example|
    if example.exception
      page = save_page
      screenshot = save_screenshot(nil)

      exception = example.exception
      exception.define_singleton_method :message do
        super() +
          "\nHTML page: #{page}\nScreenshot: #{screenshot}"
      end
    end
  end
end
