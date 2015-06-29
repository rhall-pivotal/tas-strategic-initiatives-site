require 'net/ssh/gateway'
require 'opsmgr/environments'
require 'opsmgr/cmd/bosh_command'

module Ert
  class VCloudCatsRunner
    def initialize(environment:, logger: nil)
      @environment = environment
      @logger = logger
    end

    def run_cats
      director_ip = Opsmgr::Cmd::BoshCommand.build(environment).director_ip
      gateway.open(director_ip, 25_555, 25_555) do |_|
        log_info("Opened tunnel to MicroBOSH at #{director_ip}")
        ENV['RELENG_ENV'] = environment.settings.name
        ENV['DIRECTOR_IP_OVERRIDE'] = 'localhost'
        `./scripts/runtime/post_install_test.sh`
      end
    end

    private

    attr_reader :environment, :logger

    def gateway
      uri = URI.parse(environment.settings.ops_manager.url)
      log_info("Setting up SSH gateway to OpsManager at #{uri.host}")
      Net::SSH::Gateway.new(
        uri.host,
        'ubuntu',
        password: 'tempest'
      )
    end

    def log_info(message)
      logger.info(message) if logger
    end
  end
end
