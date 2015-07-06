require 'spec_helper'
require 'ert/aws_database_creator'
require 'recursive-open-struct'

describe Ert::AwsDatabaseCreator do
  let(:settings) do
    RecursiveOpenStruct.new(
      'ops_manager' => {
        'url' => 'https://some.host.name/',
        'elastic_runtime' => {
          'rds' => {
            'host' => 'some-host',
            'port' => 555,
            'username' => 'some-user',
            'password' => 'some-password',
            'dbname' => 'some-dbname',
          },
        },
        'aws' => {
          'ssh_key' => 'ops-manager-pem-data'
        }
      }
    )
  end
  subject(:aws_database_creator) { Ert::AwsDatabaseCreator.new(settings: settings) }

  let(:mysql2_client) { instance_double(Mysql2::Client) }
  let(:gateway) { instance_double(Net::SSH::Gateway) }

  describe '#create_dbs' do
    before do
      allow(Mysql2::Client).to receive(:new).and_return(mysql2_client)
      allow(mysql2_client).to receive(:query)
      allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
      allow(gateway).to receive(:open).and_yield(444)
    end

    it 'establishes an ssh tunnel to the ops manager' do
      expect(Net::SSH::Gateway).to receive(:new).with(
        'some.host.name',
        'ubuntu',
        key_data: ['ops-manager-pem-data']
      ).and_return(gateway)

      expect(gateway).to receive(:open).with(
        settings.ops_manager.elastic_runtime.rds.host,
        settings.ops_manager.elastic_runtime.rds.port,
      ).and_yield(444)

      aws_database_creator.create_dbs
    end

    it 'connects using the db information in the settings object' do
      expect(Mysql2::Client).to receive(:new).with(
        host: '127.0.0.1',
        port: 444,
        username: settings.ops_manager.elastic_runtime.rds.username,
        password: settings.ops_manager.elastic_runtime.rds.password,
      ).and_return(mysql2_client)

      aws_database_creator.create_dbs
    end

    it 'creates dbs for all ERT apps' do
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS ccdb')
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS uaa')
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS notifications')
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS autoscale')
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS console')
      expect(mysql2_client).to receive(:query).with('CREATE DATABASE IF NOT EXISTS app_usage_service')

      aws_database_creator.create_dbs
    end
  end
end
