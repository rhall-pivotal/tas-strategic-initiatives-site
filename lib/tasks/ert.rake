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

  desc '- Run CATS in as SSH tunnel established with a password [:environment] '
  task :run_cats_ssh_password_tunnel, [:environment] do |_, args|
    require 'opsmgr/log'
    require 'opsmgr/environments'
    require 'ert/vcloud_cats_runner'
    logger = Opsmgr.logger_for('Rake')
    logger.info "Run CATS for #{args[:environment]}"

    environment = Opsmgr::Environments.for(args.environment)
    runner = Ert::VCloudCatsRunner.new(environment: environment, logger: logger)
    runner.run_cats
  end
end
