require 'aws/route_53'
require 'backport_refinements'
using OpsManagerUiDrivers::BackportRefinements

module Ert
  class DnsUpdater
    def initialize(settings:)
      env_config = settings.dig('vm_shepherd', 'env_config')
      self.name = env_config.dig('stack_name')
      self.elb_dns_name = settings.dig('ops_manager', 'elastic_runtime', 'elb_dns_name')
      self.ssh_elb_dns_name = settings.dig('ops_manager', 'elastic_runtime', 'ssh_elb_dns_name')
      self.tcp_elb_dns_name = settings.dig('ops_manager', 'elastic_runtime', 'tcp_elb_dns_name')
      self.hosted_zone_id = settings.dig('vm_shepherd', 'env_config', 'route53_hosted_zone_id')
      self.deployment_domain = settings.dig('ops_manager', 'elastic_runtime', 'deployment_domain')

      self.wildcard_fqdn = "*.#{deployment_domain}."
      self.ssh_fqdn = "ssh.#{deployment_domain}."
      self.tcp_fqdn = "tcp.#{deployment_domain}."

      self.route53 = AWS::Route53.new(
        access_key_id: env_config.dig('aws_access_key'),
        secret_access_key: env_config.dig('aws_secret_key')
      )
    end

    def update_record
      record = wildcard_record || new_a_record_set
      record[:resource_records].first[:value] = elb_dns_name
      record[:ttl] = 5

      new_ssh_record = ssh_record || new_a_record_set
      new_ssh_record[:resource_records].first[:value] = ssh_elb_dns_name
      new_ssh_record[:ttl] = 5

      new_tcp_record = tcp_record || new_a_record_set
      new_tcp_record[:resource_records].first[:value] = tcp_elb_dns_name
      new_tcp_record[:ttl] = 5

      change_record = {
        hosted_zone_id: hosted_zone_id,
        change_batch: {
          changes: [
            {
              action: 'UPSERT',
              resource_record_set: record
            },
            {
              action: 'UPSERT',
              resource_record_set: new_ssh_record
            },
            {
              action: 'UPSERT',
              resource_record_set: new_tcp_record
            }
          ]
        }
      }

      route53.client.change_resource_record_sets(change_record)
    end

    attr_accessor :route53

    private

    attr_accessor :name, :elb_dns_name, :ssh_elb_dns_name, :tcp_elb_dns_name, :hosted_zone_id, :deployment_domain

    attr_accessor :wildcard_fqdn, :ssh_fqdn, :tcp_fqdn

    def new_a_record_set
      record = Aws::Route53::Types::ResourceRecordSet.new
      record.type = "A"
      record[:resource_records] = [Aws::Route53::Types::ResourceRecord.new]
      record
    end

    def wildcard_record
      resource_record_sets = route53.client.list_resource_record_sets(hosted_zone_id: hosted_zone_id)
      record_sets = resource_record_sets[:resource_record_sets]
      record_sets.find(default_wildcard_record) do |set|
        set[:name] == @wildcard_fqdn
      end
    end

    def ssh_record
      resource_record_sets = route53.client.list_resource_record_sets(hosted_zone_id: hosted_zone_id)
      record_sets = resource_record_sets[:resource_record_sets]
      record_sets.find(default_ssh_record) do |set|
        set[:name] == @ssh_fqdn
      end
    end

    def tcp_record
      resource_record_sets = route53.client.list_resource_record_sets(hosted_zone_id: hosted_zone_id)
      record_sets = resource_record_sets[:resource_record_sets]
      record_sets.find(default_tcp_record) do |set|
        set[:name] == @tcp_fqdn
      end
    end

    def default_ssh_record
      proc do
        {
          resource_records: [
            {
              value: 'bogus'
            }
          ],
          name: @ssh_fqdn,
          type: 'CNAME',
          ttl: 5
        }
      end
    end

    def default_tcp_record
      proc do
        {
          resource_records: [
            {
              value: 'bogus'
            }
          ],
          name: @tcp_fqdn,
          type: 'CNAME',
          ttl: 5
        }
      end
    end

    def default_wildcard_record
      proc do
        {
          resource_records: [
            {
              value: 'bogus'
            }
          ],
          name: @wildcard_fqdn,
          type: 'CNAME',
          ttl: 5
        }
      end
    end
  end
end
