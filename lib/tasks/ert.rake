namespace :ert do
  require 'tools/integration_spec_runner'

  desc '-  Configure Elastic Runtime [:environment, :ert_version, :om_version]'
  task :configure, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_ert
  end
end
