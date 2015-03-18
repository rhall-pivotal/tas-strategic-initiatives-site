require 'spec_helper'
require 'cmd/runtime'
require 'opsmgr/environment'

describe Cmd::Runtime do
  before do
    allow(Opsmgr::Api::EndpointsFactory).to receive(:create).and_return(Opsmgr::Api::Version12::Endpoints.new)
  end
  let(:settings) do
    {
      environments: {
        test: {
          ops_manager_settings: {
            ip: '127.0.0.1',
            netmask: '255.255.254.0',
            gateway: '127.0.0.254',
            reserved_ip_ranges: '127.0.0.1-127.0.0.16,127.0.0.100-127.0.1.255',
            dns: '192.168.2.3',
            ntp_servers: ['time.example.com'],
            vm_password: 'test_vm_password'
          }
        }
      }
    }
  end

  let(:environment) { Opsmgr::Environment.build(:test, settings) }

  let(:endpoints) { Opsmgr::Api::Version12::Endpoints.new }
  let(:installer) { Opsmgr::Cmd::Installer.build(environment, described_class::PRODUCT_NAME) }
  let(:runtime_product) { Runtime.build(environment) }
  let(:upgrader) { Opsmgr::Cmd::Upgrader.build(environment, described_class::PRODUCT_NAME) }
  subject(:runtime) { Cmd::Runtime.new(installer, upgrader, runtime_product) }

  describe '.build' do
    it 'correctly builds a Runtime command' do
      expect(Cmd::Runtime.build(environment)).to be_a(Cmd::Runtime)
    end
  end

  describe '#upgrade' do
    it 'calls Opsmgr::Cmd::Upgrader#upgrade and then Installer#install' do
      expect(upgrader).to receive(:upgrade).with(no_args).ordered
      expect(runtime_product).to receive(:configure).with(no_args).ordered
      expect(installer).to receive(:install).with(no_args).ordered

      runtime.upgrade
    end
  end
end
