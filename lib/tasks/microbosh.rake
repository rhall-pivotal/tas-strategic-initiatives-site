namespace :microbosh do
  require 'tools/integration_spec_runner'

  desc '-  Configure Microbosh [:environment, :om_version]'
  task :configure, [:environment, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      om_version: args.om_version
    ).configure_microbosh
  end

  desc '-  Install Microbosh [:environment, :om_version]'
  task :install, [:environment, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      om_version: args.om_version
    ).install_microbosh
  end
end
