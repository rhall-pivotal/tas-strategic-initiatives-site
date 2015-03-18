require 'opsmgr/cmd/installer'
require 'opsmgr/cmd/upgrader'
require 'runtime'

module Cmd
  class Runtime
    PRODUCT_NAME = 'cf'.freeze

    def self.build(environment)
      installer = Opsmgr::Cmd::Installer.build(environment, PRODUCT_NAME)
      upgrader = Opsmgr::Cmd::Upgrader.build(environment, PRODUCT_NAME)
      runtime_product = ::Runtime.build(environment)
      new(installer, upgrader, runtime_product)
    end

    def initialize(installer, upgrader, runtime_product)
      @installer = installer
      @upgrader = upgrader
      @runtime_product = runtime_product
    end

    def upgrade
      @upgrader.upgrade
      @runtime_product.configure
      @installer.install
    end
  end
end
