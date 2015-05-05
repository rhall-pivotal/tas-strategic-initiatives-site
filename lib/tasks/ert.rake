namespace :ert do
  require 'tools/integration_spec_runner'

  desc '-  Upload/Add Elastic Runtime [:environment, :om_version, :product_path]'
  task :upload, [:environment, :om_version, :product_path] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      om_version: args.om_version
    ).upload_ert(args.product_path)
  end

  desc '-  Configure Elastic Runtime [:environment, :ert_version, :om_version]'
  task :configure, [:environment, :ert_version, :om_version] do |_, args|
    IntegrationSpecRunner.new(
      environment: args.environment,
      ert_version: args.ert_version,
      om_version: args.om_version
    ).configure_ert
  end
end
