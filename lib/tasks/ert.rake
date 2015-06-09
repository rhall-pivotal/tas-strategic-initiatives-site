namespace :ert do
  require 'tools/integration_spec_runner'

  desc 'Create AWS Databases for Elastic Runtime [:environment]'
  task :create_aws_dbs, [:environment] do |_, args|
    require 'opsmgr/log'
    require 'opsmgr/environments'
    require 'ert/aws_database_creator'

    logger = Opsmgr.logger_for('Rake')
    logger.info "Creating AWS DBs for #{args[:environment]}"

    environment = Opsmgr::Environments.for(args.environment)
    creator = Ert::AwsDatabaseCreator.new(settings: environment.settings)
    creator.create_dbs
  end

  desc 'Configure Elastic Runtime [:environment, :ert_version, :om_version]'
  task :configure, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_ert
  end
end
