require 'opsmgr/environments'
require 'opsmgr/cmd/bosh_command'
require 'open3'
require 'net/ssh/gateway'

module Ert
  class CatsRunner
    def initialize(environment_name:, logger:)
      @environment = Opsmgr::Environments.for(environment_name)
      @logger = logger
    end

    def run_cats
      gateway do |_|
        set_bosh_deployment

        system_or_fail(
          "#{bosh_command_prefix} run errand #{errand_name}",
          'CF Acceptance Tests failed'
        )
      end
    end

    private

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

    def set_bosh_deployment
      system_or_fail(bosh_command.target, 'bosh target failed')

      bosh_deployment = bosh_deployment_name(bosh_command_prefix)

      deployment_file = "#{ENV.fetch('TMPDIR', '/tmp')}/#{environment.settings.name}.yml"

      system_or_fail(
        "#{bosh_command_prefix} -n download manifest #{bosh_deployment} #{deployment_file}",
        'bosh download manifest failed'
      )
      system_or_fail(
        "#{bosh_command_prefix} deployment #{deployment_file}",
        'bosh deployment failed'
      )
    end

    def system_or_fail(command, failure_message)
      logger.info("Running #{command}")
      Bundler.clean_system(command) || fail(failure_message)
    end

    def bosh_deployment_name(command)
      @bosh_deployment_name ||= begin
        Bundler.with_clean_env do
          bosh_deployment, status = Open3.capture2("#{command} deployments | grep -Eoh 'cf-[0-9a-f]{8,}'")
          fail('bosh deployments failed') unless status.success?
          bosh_deployment.chomp
        end
      end
    end

    def bosh_command
      @bosh_command ||= Opsmgr::Cmd::BoshCommand.build(environment)
    end

    def bosh_command_prefix
      @bosh_command_prefix ||= bosh_command.command
    end

    def errand_name
      environment.settings.internetless ? 'acceptance-tests-internetless' : 'acceptance-tests'
    end

    def ssh_key
      case environment.settings.iaas_type
      when 'aws'
        environment.settings.ops_manager.aws.ssh_key
      when 'openstack'
        environment.settings.ops_manager.openstack.ssh_private_key
      end
    end

    attr_reader :environment, :logger
  end
end
