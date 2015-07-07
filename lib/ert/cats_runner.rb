require 'opsmgr/environments'
require 'opsmgr/cmd/bosh_command'
require 'open3'

module Ert
  class CatsRunner
    def initialize(environment_name:)
      @environment = Opsmgr::Environments.for(environment_name)
    end

    def run_cats
      set_bosh_deployment

      system_or_fail(
        "#{bosh_command_prefix} run errand #{errand_name}",
        'CF Acceptance Tests failed'
      )
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

    attr_reader :environment
  end
end
