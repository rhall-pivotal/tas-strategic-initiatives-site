namespace :ert do
  require 'tools/integration_spec_runner'

  desc 'Create AWS Databases for Elastic Runtime'
  task :create_aws_dbs, [:environment_name] do |_, args|
    require 'opsmgr/log'
    require 'opsmgr/environments'
    require 'ert/aws_database_creator'

    logger = Opsmgr.logger_for('Rake')
    environment = Opsmgr::Environments.for(args.environment_name)

    iaas = environment.settings.iaas_type
    if (iaas == 'aws')
      logger.info "Creating AWS DBs for #{args[:environment_name]}"

      creator = Ert::AwsDatabaseCreator.new(settings: environment.settings)
      creator.create_dbs
    else
      logger.info "Not creating AWS databases because environment is: #{iaas}"
    end
  end

  desc 'Update DNS for an ELB for Elastic Runtime'
  task :update_dns_elb, [:environment_name] do |_, args|
    require 'opsmgr/log'
    require 'opsmgr/environments'
    require 'ert/dns_updater'

    logger = Opsmgr.logger_for('Rake')

    environment = Opsmgr::Environments.for(args.environment_name)
    iaas = environment.settings.iaas_type
    if (iaas == 'aws')
      logger.info "Updating DNS record for #{args[:environment_name]}"

      dns_updater = Ert::DnsUpdater.new(settings: environment.settings)
      dns_updater.update_record
    else
      logger.info "Not updating ELB's DNS because environment is: #{iaas}"
    end
  end

  desc 'Configure Elastic Runtime'
  task :configure, [:environment_name, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment_name,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_ert
  end

  desc 'Configure Elastic Runtime External Databases'
  task :configure_external_dbs, [:environment_name, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment_name,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_external_dbs
  end

  desc 'Configure Elastic Runtime External File Storage'
  task :configure_external_file_storage, [:environment_name, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment_name,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_external_file_storage
  end

  desc 'run the cats errand'
  task :run_cats, [:environment_name] do |_, args|
    require 'ert/cats_runner'
    require 'opsmgr/log'

    logger = Opsmgr.logger_for('Rake')
    Ert::CatsRunner.new(environment_name: args.environment_name, logger: logger).run_cats
  end

  namespace :pipeline do
    desc 'create a feature pipeline'
    task :create_feature_pipeline, [:branch_name, :iaas_type] do |_, args|
      require 'pipeline/feature_pipeline_creator'
      Pipeline::FeaturePipelineCreator.new(
        branch_name: args.branch_name,
        iaas_type: args.iaas_type
      ).create_pipeline
    end

    desc 'create a feature upgrade pipeline'
    task :create_upgrade_pipeline, [:branch_name, :iaas_type, :ert_initial_full_version, :om_initial_full_version] do |_, args|
      require 'pipeline/feature_pipeline_creator'
      Pipeline::FeaturePipelineCreator.new(
        branch_name: args.branch_name,
        iaas_type: args.iaas_type
      ).create_upgrade_pipeline(
        ert_initial_full_version: args.ert_initial_full_version,
        om_initial_full_version: args.om_initial_full_version
      )
    end

    desc 'deploy feature pipeline to concourse'
    task :deploy_feature_pipeline, [:branch_name] do |_, args|
      require 'pipeline/feature_pipeline_deployer'
      Pipeline::FeaturePipelineDeployer.new(
        branch_name: args.branch_name,
      ).deploy_pipeline
    end
  end
end
