require 'spec_helper'
require 'ert/vcloud_cats_runner'
require 'recursive-open-struct'

describe Ert::VCloudCatsRunner do
  subject(:vcloud_cats_runner) { Ert::VCloudCatsRunner.new(environment: environment) }

  let(:gateway) { instance_double(Net::SSH::Gateway) }
  let(:environment) { instance_double(Opsmgr::Environments) }
  let(:bosh_command) { instance_double(Opsmgr::Cmd::BoshCommand) }

  let(:microbosh_ip_address) { 'some.ip.address' }
  let(:microbosh_port) { 25_555 }
  let(:local_port) { 25_555 }

  let(:settings) do
    RecursiveOpenStruct.new(
      'name' => 'some-environment',
      'ops_manager' => {
        'url' => 'https://some.host.name',
      }
    )
  end

  before do
    allow(environment).to receive(:settings).and_return(settings)
    allow(Opsmgr::Cmd::BoshCommand).to receive(:build).with(environment).and_return(bosh_command)
    allow(bosh_command).to receive(:director_ip).and_return(microbosh_ip_address)
    allow(vcloud_cats_runner).to receive(:'`')
  end

  describe '#run_cats' do
    before do
      allow(Net::SSH::Gateway).to receive(:new).and_return(gateway)
      allow(gateway).to receive(:open).and_yield(444)
    end

    it 'establishes an ssh tunnel to the ops manager' do
      expect(Net::SSH::Gateway).to receive(:new).with(
        'some.host.name',
        'ubuntu',
        password: 'tempest'
      ).and_return(gateway)

      expect(gateway).to receive(:open).with(
        microbosh_ip_address,
        microbosh_port,
        local_port,
      ).and_yield(444)

      vcloud_cats_runner.run_cats
    end

    it 'shells out to post_install_test.sh script' do
      expect(ENV).to receive(:[]=).with('RELENG_ENV', 'some-environment')
      expect(ENV).to receive(:[]=).with('DIRECTOR_IP_OVERRIDE', 'localhost')
      expect(vcloud_cats_runner).to receive(:'`').with('./scripts/runtime/post_install_test.sh')
      vcloud_cats_runner.run_cats
    end
  end
end
