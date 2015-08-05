require 'net/ssh/gateway'
require 'opsmgr/cmd/bosh_command'

module Ert
  class IaasGateway
    def initialize(bosh_command:, environment_name:, logger:)
      @bosh_command = bosh_command
      @environment = Opsmgr::Environments.for(environment_name)
      @logger = logger
    end

    def gateway(&block)
      case environment.settings.iaas_type
      when 'vsphere'
        block.call
      when 'aws', 'openstack'
        ssh_key_gateway(block)
      when 'vcloud'
        ssh_password_gateway(block)
      end
    end

    private

    attr_reader :bosh_command, :environment, :logger

    def ssh_password_gateway(block)
      ENV['DIRECTOR_IP_OVERRIDE'] = 'localhost'
      uri = URI.parse(environment.settings.ops_manager.url)
      director_ip = bosh_command.director_ip
      logger.info("Setting up SSH gateway to OpsManager at #{uri.host}")
      Net::SSH::Gateway.new(
        uri.host,
        'ubuntu',
        password: 'tempest'
      ).open(director_ip, 25_555, 25_555) do |_|
        logger.info("Opened tunnel to MicroBOSH at #{director_ip}")
        block.call
      end
    end

    def ssh_key_gateway(block)
      ENV['DIRECTOR_IP_OVERRIDE'] = 'localhost'
      uri = URI.parse(environment.settings.ops_manager.url)
      director_ip = bosh_command.director_ip
      Net::SSH::Gateway.new(
        uri.host,
        'ubuntu',
        key_data: [ssh_key]
      ).open(director_ip, 25_555, 25_555) do |_|
        logger.info("Opened tunnel to MicroBOSH at #{director_ip}")
        block.call
      end
    end

    def ssh_key
      case environment.settings.iaas_type
      when 'aws'
        environment.settings.ops_manager.aws.ssh_key
      when 'openstack'
        environment.settings.ops_manager.openstack.ssh_private_key
      end
    end
  end
end
