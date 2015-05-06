require 'rspec'

class IntegrationSpecRunner
  class UnsupportedErtVersion < StandardError
  end

  INTEGRATION_SPEC_PREFIX = 'integration'
  SUPPORTED_ERT_VERSIONS = %w(1.4 1.5)

  def initialize(environment:, om_version:, ert_version: nil)
    ENV['ENVIRONMENT_NAME'] = environment
    ENV['OM_VERSION'] = om_version
    if ert_version.nil? || SUPPORTED_ERT_VERSIONS.include?(ert_version)
      @ert_version = ert_version
    else
      fail UnsupportedErtVersion, "Version #{ert_version.inspect} is not supported"
    end
  end

  def upload_ert(product_path)
    ENV['PRODUCT_PATH'] = product_path
    RSpec::Core::Runner.run([ert_spec_path + '/upload_ert_spec.rb'])
  end

  def configure_ert
    RSpec::Core::Runner.run([ert_spec_path + '/configure_ert_spec.rb'])
  end

  def configure_microbosh
    RSpec::Core::Runner.run([microbosh_spec_path + '/configure_microbosh_spec.rb'])
  end

  def install_microbosh
    RSpec::Core::Runner.run([microbosh_spec_path + '/install_microbosh_spec.rb'])
  end

  private

  attr_reader :ert_version

  def ert_spec_path
    "#{INTEGRATION_SPEC_PREFIX}/ERT-#{ert_version}"
  end

  def microbosh_spec_path
    "#{INTEGRATION_SPEC_PREFIX}/microbosh"
  end
end
