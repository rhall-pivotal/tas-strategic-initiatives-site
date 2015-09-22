require 'opsmgr/environments'
require 'open3'

module Ert
  class NotificationsTestsRunner
    def initialize(iaas_gateway:, bosh_command:, environment_name:, logger:)
      @iaas_gateway = iaas_gateway
      @bosh_command = bosh_command
      @environment = Opsmgr::Environments.for(environment_name)
      @logger = logger
    end

    def run
      iaas_gateway.gateway do
        set_bosh_deployment

        system_or_fail(
          "#{bosh_command_prefix} run errand notifications-tests",
          'CF Acceptance Tests failed'
        )

        system_or_fail(
          "#{bosh_command_prefix} run errand notifications-ui-tests",
          'CF Acceptance Tests failed'
        )
      end
    end

    private

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

    def bosh_command_prefix
      @bosh_command_prefix ||= bosh_command.command
    end

    attr_reader :iaas_gateway, :bosh_command, :environment, :logger
  end
end
