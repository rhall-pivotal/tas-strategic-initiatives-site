require 'spec_helper'
require 'yaml'

describe 'handcraft' do
  describe 'syslog_aggregator_config' do
    handcraft = YAML.load(File.read(File.join(SPEC_ROOT, '..', 'metadata_parts', 'handcraft.yml')))

    jobs_with_metron = handcraft['job_types'].select do |job_type|
      next unless job_type['templates']

      job_type['templates'].map { |template| template['name'] }.include?('metron_agent')
    end

    jobs_with_metron.each do |job|
      it "is present in #{job['name']} job" do
        manifest_snippet = YAML.load(job['manifest'])
        expect(manifest_snippet).to include('syslog_daemon_config')
      end
    end
  end
end
