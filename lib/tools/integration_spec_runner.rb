require 'rspec'

class IntegrationSpecRunner
  class UnsupportedErtVersion < StandardError
  end

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

  def configure_ert
    RSpec::Core::Runner.run(["integration/ERT-#{ert_version}/configure_ert_spec.rb"])
  end

  private

  attr_reader :ert_version
end
