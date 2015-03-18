desc 'Deploy Elastic Runtime'

namespace :runtime do
  desc '-  Upload Elastic Runtime product (cf-*.pivotal)'
  task :upload, [:environment, :product_path] do |_, args|
    require 'opsmgr/cmd/uploader'
    require 'opsmgr/environments'

    Opsmgr::Cmd::Uploader.build(Opsmgr::Environments.for(args.environment), File.expand_path(args.product_path)).upload
  end

  desc '-  Install Elastic Runtime'
  task :install, [:environment] do |_, args|
    require 'opsmgr/cmd/installer'
    require 'opsmgr/environments'

    Opsmgr::Cmd::Installer.build(Opsmgr::Environments.for(args.environment), 'Elastic Runtime').install
  end
end
