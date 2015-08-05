require 'spec_helper'
require 'ert/internet_checker'

describe Ert::InternetChecker do
  subject(:internet_checker) do
    Ert::InternetChecker.new(environment_name: environment_name)
  end

  let(:environment_name) { 'env_name' }

  let(:env_config) do
    <<YAML
---
name: env_name
iaas_type: 'vsphere'
ops_manager:
  url: https://pcf.foo.cf-app.com
YAML
  end

  let(:settings) do
    RecursiveOpenStruct.new(YAML.load(env_config), recurse_over_arrays: true)
  end

  let(:opsmgr_environment) { instance_double(Opsmgr::Environments, settings: settings) }

  before do
    allow(Opsmgr::Environments).to receive(:for).and_return(opsmgr_environment)
  end

  describe '#internetless?' do
    let(:google_80) { false }
    let(:google_443) { false }
    let(:google_25) { true }
    let(:google_587) { true }
    let(:github_22) { false }
    before do
      allow(subject).to receive(:connection_allowed?).with('www.google.com', 80).and_return(google_80)
      allow(subject).to receive(:connection_allowed?).with('www.google.com', 443).and_return(google_443)
      allow(subject).to receive(:connection_allowed?).with('smtp-relay.gmail.com', 25).and_return(google_25)
      allow(subject).to receive(:connection_allowed?).with('smtp-relay.gmail.com', 587).and_return(google_587)
      allow(subject).to receive(:connection_allowed?).with('github.com', 22).and_return(github_22)
    end

    it 'returns true if all conditions are met' do
      expect(subject).to receive(:connection_allowed?).with('www.google.com', 80)
      expect(subject).to receive(:connection_allowed?).with('www.google.com', 443)
      expect(subject).to receive(:connection_allowed?).with('smtp-relay.gmail.com', 25)
      expect(subject).to receive(:connection_allowed?).with('smtp-relay.gmail.com', 587)
      expect(subject).to receive(:connection_allowed?).with('github.com', 22)

      expect(subject.internetless?).to be_truthy
    end

    context 'can connect to www.google.com:80' do
      let(:google_80) { true }

      it 'returns false' do
        expect(subject.internetless?).to be_falsey
      end
    end

    context 'can connect to www.google.com:443' do
      let(:google_443) { true }

      it 'returns false' do
        expect(subject.internetless?).to be_falsey
      end
    end

    context 'can not connect to smtp-relay.gmail.com:80' do
      let(:google_25) { false }

      it 'returns false' do
        expect(subject.internetless?).to be_falsey
      end
    end

    context 'can not connect to smtp-relay.gmail.com:80' do
      let(:google_587) { false }

      it 'returns false' do
        expect(subject.internetless?).to be_falsey
      end
    end

    context 'can connect to github.com:22' do
      let(:github_22) { true }

      it 'returns false' do
        expect(subject.internetless?).to be_falsey
      end
    end
  end

  describe '#connection_allowed?' do
    let(:host) { 'hostname' }
    let(:port) { 'port_number' }
    let(:ssh) { double('ssh session') }
    let(:result) { ['', '', 0, ''] }
    let(:ops_manager_hostname) { 'pcf.foo.cf-app.com' }
    let(:ops_manager_username) { 'ubuntu' }
    let(:ops_manager_password) { 'tempest' }

    before do
      allow(Net::SSH).to receive(:start).and_yield(ssh)
      allow(subject).to receive(:ssh_exec!).and_return(result)
    end

    it 'creates a ssh connection using the right hostname, username, and password' do
      expect(Net::SSH).to(
        receive(:start).with(
          ops_manager_hostname,
          ops_manager_username,
          password: ops_manager_password
        ).and_yield(ssh)
      )
      subject.connection_allowed?(host, port)
    end

    it 'runs the nc command with the right host and port' do
      expect(subject).to receive(:ssh_exec!).with(ssh, "echo QUIT | nc -w 5 #{host} #{port}")
      subject.connection_allowed?(host, port)
    end

    context 'when the connection attempt fails' do
      let(:result) { ['', '', 1, ''] }

      it 'returns false' do
        expect(subject.connection_allowed?(host, port)).to eq(false)
      end
    end

    context 'when the connection attempt succeeds' do
      let(:result) { ['', '', 0, ''] }

      it 'returns true' do
        expect(subject.connection_allowed?(host, port)).to eq(true)
      end
    end
  end
end
