require 'spec_helper'
require 'recursive-open-struct'
require 'ert/dns_updater'
require 'aws/route_53'

describe Ert::DnsUpdater do
  let(:settings) do
    {
      'name' => 'some',
      'vm_shepherd' => {
        'env_config' => {
          'aws_access_key' => 'some-access-key',
          'aws_secret_key' => 'some-secret-key'
        }
      },
      'ops_manager' => {
        'elastic_runtime' => {
          'elb_dns_name' => 'some-fake-elb-dns-name',
          'ssh_elb_dns_name' => 'some-fake-ssh-elb-dns-name'
        }
      }
    }
  end
  let(:hosted_zone_id) { 'some-hosted-zone' }
  let(:hosted_zones) do
    {
      hosted_zones: [
        {
          id: hosted_zone_id,
          name: 'some.hosted.name',
          caller_reference: 'some-caller-reference',
          config: { private_zone: false },
          resource_record_set_count: 4
        },
        {
          id: 'another-hosted-zone',
          name: 'another.hosted.name',
          caller_reference: 'another-caller-reference',
          config: { private_zone: false },
          resource_record_set_count: 4
        },
      ]
    }
  end
  let(:record_sets) do
    {
      resource_record_sets: [
        {
          resource_records: [
            { value: 'fake-hostname' },
            { value: 'another-fake-hostname' },
          ],
          name: 'some.cf-app.com.',
          type: 'NS', ttl: 172_800
        },
        {
          resource_records: [
            {
              value: 'another-crazy-fake-dns'
            }
          ],
          name: 'some.cf-app.com.',
          type: 'SOA', ttl: 900
        },
        wildcard_record,
        ssh_record
      ].compact
    }
  end
  let(:wildcard_record) do
    {
      resource_records: [
        { value: 'elb-pcf-specific-hostname' }
      ],
      name: '\\052.some.cf-app.com.',
      type: 'CNAME', ttl: 300
    }
  end
  let(:ssh_record) do
    {
      resource_records: [
        { value: 'ssh-elb-pcf-specific-hostname' }
      ],
      name: 'ssh.some.cf-app.com.',
      type: 'CNAME', ttl: 300
    }
  end

  let(:r53) { instance_double(AWS::Route53) }
  let(:client) { instance_double(AWS::Route53::Client::V20130401) }

  subject(:dns_updater) { Ert::DnsUpdater.new(settings: settings) }

  before do
    allow(AWS::Route53).to receive(:new)
      .with(
        access_key_id: 'some-access-key',
        secret_access_key: 'some-secret-key')
      .and_return(r53)
    allow(r53).to receive(:client).and_return(client)
    allow(client).to receive(:list_hosted_zones).and_return(hosted_zones)
    allow(client).to receive(:list_resource_record_sets).and_return(record_sets)
  end

  describe '#initialize' do
    it 'creates a route53 object using environment variables for credentials' do
      expect(dns_updater.route53).to eq(r53)
    end
  end

  describe '#update_record' do
    let(:updated_record_set) do
      {
        hosted_zone_id: hosted_zone_id,
        change_batch: {
          changes: [
            {
              action: 'UPSERT',
              resource_record_set: {
                resource_records: [
                  { value: 'some-fake-elb-dns-name' }
                ],
                name: '\\052.some.cf-app.com.',
                type: 'CNAME',
                ttl: 5
              }

            },
            {
              action: 'UPSERT',
              resource_record_set: {
                resource_records: [
                  { value: 'some-fake-ssh-elb-dns-name' }
                ],
                name: 'ssh.some.cf-app.com.',
                type: 'CNAME',
                ttl: 5
              }

            }
          ]
        }
      }
    end
    let(:data_response) do
      {
        change_info: {
          id: 'some-id',
          status: 'some-status',
        }
      }
    end
    let(:update_response) { instance_double(AWS::Core::Response) }

    before do
      allow(client).to receive(:change_resource_record_sets).with(updated_record_set)
        .and_return(data_response)
    end

    it 'updates the wildcard DNS record with the ELB information' do
      expect(dns_updater.update_record).to eq(data_response)
    end

    context 'when there is no existing ssh record' do
      let(:ssh_record) { nil }
      it 'uses the default record as a template' do
        expect(client).to receive(:change_resource_record_sets).with(updated_record_set)
          .and_return(data_response)
        expect(dns_updater.update_record).to eq(data_response)
      end
    end

    context 'when there is no existing wildcard record' do
      let(:wildcard_record) { nil }
      it 'uses the default record as a template' do
        expect(client).to receive(:change_resource_record_sets).with(updated_record_set)
          .and_return(data_response)
        expect(dns_updater.update_record).to eq(data_response)
      end
    end
  end
end
