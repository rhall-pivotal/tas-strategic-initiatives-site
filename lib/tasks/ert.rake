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

  desc 'Update DNS for an ELB for Elastic Runtime [:environment]'
  task :update_dns_elb, [:environment] do |_, args|
    require 'opsmgr/log'
    require 'opsmgr/environments'
    require 'ert/dns_updater'

    logger = Opsmgr.logger_for('Rake')
    logger.info "Updating DNS record for #{args[:environment]}"

    environment = Opsmgr::Environments.for(args.environment)
    dns_updater = Ert::DnsUpdater.new(settings: environment.settings)
    dns_updater.update_record
  end

  desc 'Configure Elastic Runtime [:environment, :ert_version, :om_version]'
  task :configure, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_ert
  end

  desc 'Configure Elastic Runtime External Databases [:environment, :ert_version, :om_version]'
  task :configure_external_dbs, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_external_dbs
  end

  desc 'Configure Elastic Runtime External File Storage [:environment, :ert_version, :om_version]'
  task :configure_external_file_storage, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_external_file_storage
  end
  namespace :pipeline do
    desc 'create a feature pipeline'
    task :create_feature_pipeline, [:branch_name, :iaas_type] do |_, args|
      require 'pipeline/create_feature_pipeline'
      Pipeline::CreateFeaturePipeline.new(
        branch_name: args.branch_name,
        iaas_type: args.iaas_type
      ).create_pipeline
    end
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
