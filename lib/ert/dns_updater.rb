require 'aws/route_53'

module Ert
  class DnsUpdater
    def initialize(settings:)
      self.name = settings.name
      self.elb_dns_name = settings.ops_manager.elastic_runtime.elb_dns_name
      access_key = settings.vm_shepherd.env_config.aws_access_key
      secret_access_key = settings.vm_shepherd.env_config.aws_secret_key
      self.route53 = AWS::Route53.new(
        access_key_id: access_key,
        secret_access_key: secret_access_key
      )
    end

    def update_record
      record = wildcard_record
      record[:resource_records].first[:value] = elb_dns_name
      record[:ttl] = 30
      change_record = {
        hosted_zone_id: hosted_zone_id,
        change_batch: {
          changes: [{
            action: 'UPSERT',
            resource_record_set: record
          }]
        }
      }
      route53.client.change_resource_record_sets(change_record)
    end

    attr_accessor :route53

    private

    attr_accessor :name, :elb_dns_name

    def hosted_zone_id
      resp = route53.client.list_hosted_zones
      resp[:hosted_zones].find do |zone|
        zone[:name].include? name
      end[:id]
    end

    def wildcard_record
      resource_record_sets = route53.client.list_resource_record_sets(hosted_zone_id: hosted_zone_id)
      record_sets = resource_record_sets[:resource_record_sets]
      record_sets.find do |set|
        set[:name].include? '052'
      end
    end
  end
end
