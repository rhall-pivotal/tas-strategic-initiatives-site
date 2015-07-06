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
    run_spec(["integration/ERT-#{ert_version}/configure_ert_spec.rb"])
  end

  def configure_external_dbs
    run_spec(["integration/ERT-#{ert_version}/configure_external_dbs_spec.rb"])
  end

  def configure_external_file_storage
    run_spec(["integration/ERT-#{ert_version}/configure_external_file_storage_spec.rb"])
  end

  private

  def run_spec(spec_to_run)
    RSpecExiter.exit_rspec(RSpec::Core::Runner.run(spec_to_run))
  end

  attr_reader :ert_version
end

module RSpecExiter
  def self.exit_rspec(exit_code)
    exit exit_code
  end
end
