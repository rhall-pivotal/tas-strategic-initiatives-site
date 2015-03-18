desc 'Deploy Elastic Runtime'
task :runtime,
     [:environment, :product_path] => %w(runtime:provision)

namespace :runtime do
  desc '- Alias for :upload and :configure'
  task :setup,
       [:environment, :product_path] =>
         %w(runtime:upload runtime:configure)

  desc '- Alias for :upload, :configure, and :install'
  task :provision,
       [:environment, :product_path] =>
         %w(runtime:setup runtime:install)

  desc '-  Upload Elastic Runtime product (cf-*.pivotal)'
  task :upload, [:environment, :product_path] do |_, args|
    require 'opsmgr/cmd/uploader'
    require 'opsmgr/environments'

    Opsmgr::Cmd::Uploader.build(Opsmgr::Environments.for(args.environment), File.expand_path(args.product_path)).upload
  end

  desc '-  Configure Elastic Runtime'
  task :configure, [:environment] do |_, args|
    require 'opsmgr/configurator'
    require 'runtime'
    require 'opsmgr/environments'

    Runtime.new(
      Opsmgr::Configurator.build(Opsmgr::Environments.for(args.environment))
    ).configure
  end

  desc '-  Install Elastic Runtime'
  task :install, [:environment] do |_, args|
    require 'opsmgr/cmd/installer'
    require 'opsmgr/environments'

    Opsmgr::Cmd::Installer.build(Opsmgr::Environments.for(args.environment), 'Elastic Runtime').install
  end
end
