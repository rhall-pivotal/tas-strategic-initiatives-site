require 'mysql2'
require 'net/ssh/gateway'
require 'backport_refinements'
using OpsManagerUiDrivers::BackportRefinements

module Ert
  class AwsDatabaseCreator
    def initialize(settings:)
      self.settings = settings
    end

    def create_dbs
      rds_settings = settings.dig('ops_manager', 'elastic_runtime', 'rds')
      gateway.open(rds_settings.dig('host'), rds_settings.dig('port')) do |port|
        mysql_client = Mysql2::Client.new(
          host: '127.0.0.1',
          port: port,
          username: rds_settings.dig('username'),
          password: rds_settings.dig('password')
        )

        %w(ccdb uaa notifications autoscale console app_usage_service).each do |db_name|
          mysql_client.query("CREATE DATABASE IF NOT EXISTS #{db_name}")
        end
      end
    end

    private

    attr_accessor :settings

    def gateway
      uri = URI.parse(settings.dig('ops_manager', 'url'))
      Net::SSH::Gateway.new(
        uri.host,
        'ubuntu',
        key_data: [settings.dig('ops_manager', 'aws', 'ssh_key')]
      )
    end
  end
end
