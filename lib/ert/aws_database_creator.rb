require 'mysql2'
require 'net/ssh/gateway'

module Ert
  class AwsDatabaseCreator
    def initialize(settings:)
      self.settings = settings
    end

    def create_dbs
      mysql_settings = settings.ops_manager.mysql
      gateway.open(mysql_settings.host, mysql_settings.port) do |port|
        mysql_client = Mysql2::Client.new(
          host: '127.0.0.1',
          port: port,
          username: mysql_settings.user,
          password: mysql_settings.password
        )

        %w(ccdb uaa notifications autoscale console app_usage_service).each do |db_name|
          mysql_client.query("create database #{db_name}")
        end
      end
    end

    private

    attr_accessor :settings

    def gateway
      uri = URI.parse(settings.ops_manager.url)
      Net::SSH::Gateway.new(
        uri.host,
        'ubuntu',
        key_data: [settings.ops_manager.aws.ssh_key]
      )
    end
  end
end
