require 'spec_helper'
require 'ert/iaas_gateway'

describe Ert::IaasGateway do
  subject(:iaas_gateway) do
    Ert::IaasGateway.new(
      bosh_command: bosh_command,
      environment_name: environment_name,
      logger: logger
    )
  end

  let(:environment_name) { 'env_name' }

  let(:env_config) do
    <<YAML
---
name: env_name
iaas_type: 'vsphere'
YAML
  end

  let(:settings) do
    RecursiveOpenStruct.new(YAML.load(env_config), recurse_over_arrays: true)
  end

  let(:opsmgr_environment) { instance_double(Opsmgr::Environments, settings: settings) }
  let(:bosh_command) { instance_double(Opsmgr::Cmd::BoshCommand) }
  let(:deployment_name) { 'cf-deadbeef12345678' }
  let(:logger) { instance_double(Opsmgr::LoggerWithProgName) }
  let(:gateway) { instance_double(Net::SSH::Gateway) }

  before do
    allow(Opsmgr::Environments).to receive(:for).and_return(opsmgr_environment)
    allow(logger).to receive(:info)
  end

  it 'uses the config for the given environment' do
    expect(Opsmgr::Environments).to receive(:for).with(environment_name)
    iaas_gateway
  end

  context 'vsphere' do
    it 'does not set up the ssh gateway' do
      expect(Net::SSH::Gateway).to_not receive(:new)
      expect(gateway).to_not receive(:open)
      expect(bosh_command).to_not receive(:director_ip)

      iaas_gateway.gateway {}
    end

    it 'does not set the DIRECTOR_IP_OVERRIDE environment var' do
      expect(ENV).to_not receive(:[]=)
      iaas_gateway.gateway {}
    end

    it 'yields to the block' do
      expect { |b| iaas_gateway.gateway(&b) }.to yield_control.exactly(1).times
    end
  end

  context 'aws' do
    let(:env_config) do
      <<YAML
---
name: env_name
iaas_type: aws
ops_manager:
  url: https://foo.com
  aws:
    ssh_key: key
YAML
    end

    let(:bosh_command) { instance_double(Opsmgr::Cmd::BoshCommand) }

    before do
      allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
      allow(gateway).to receive(:open).and_yield(25_555)
      allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
    end

    it 'sets up the ssh gateway' do
      expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', key_data: ['key']).and_return(gateway)
      expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
      expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

      iaas_gateway.gateway {}
    end

    it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
      expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
      iaas_gateway.gateway {}
    end

    it 'yields to the block' do
      expect { |b| iaas_gateway.gateway(&b) }.to yield_control.exactly(1).times
    end
  end

  context 'vcloud' do
    let(:env_config) do
      <<YAML
---
name: env_name
iaas_type: vcloud
ops_manager:
  url: https://foo.com
YAML
    end

    before do
      allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
      allow(gateway).to receive(:open).and_yield(25_555)
      allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
    end

    it 'sets up the ssh gateway' do
      expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', password: 'tempest').and_return(gateway)
      expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
      expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

      iaas_gateway.gateway {}
    end

    it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
      expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
      iaas_gateway.gateway {}
    end

    it 'yields to the block' do
      expect { |b| iaas_gateway.gateway(&b) }.to yield_control.exactly(1).times
    end
  end

  context 'openstack' do
    let(:env_config) do
      <<YAML
---
name: env_name
iaas_type: openstack
ops_manager:
  url: https://foo.com
  openstack:
    ssh_private_key: key
YAML
    end

    before do
      allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
      allow(gateway).to receive(:open).and_yield(25_555)
      allow(bosh_command).to receive(:director_ip).and_return('1.2.3.4')
    end

    it 'sets up the ssh gateway' do
      expect(Net::SSH::Gateway).to receive(:new).with('foo.com', 'ubuntu', key_data: ['key']).and_return(gateway)
      expect(gateway).to receive(:open).with('1.2.3.4', 25_555, 25_555).and_yield(25_555)
      expect(bosh_command).to receive(:director_ip).and_return('1.2.3.4')

      iaas_gateway.gateway {}
    end

    it 'sets the DIRECTOR_IP_OVERRIDE environment var' do
      expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
      iaas_gateway.gateway {}
    end

    it 'yields to the block' do
      expect { |b| iaas_gateway.gateway(&b) }.to yield_control.exactly(1).times
    end
  end
end
