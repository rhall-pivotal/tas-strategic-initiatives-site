require 'pipeline/suite_pipeline_creator.rb'

module Pipeline
  class FullSuitePipelineCreator < Pipeline::SuitePipelineCreator
    PIPELINES = [
      {
        method: :clean_pipeline_jobs,
        params: {
          pipeline_name: 'aws-clean',
          iaas_type: 'aws'
        },
      },
      {
        method: :clean_pipeline_jobs,
        params: {
          pipeline_name: 'openstack-clean',
          iaas_type: 'openstack'
        },
      },
      {
        method: :clean_pipeline_jobs,
        params: {
          pipeline_name: 'vsphere-clean',
          iaas_type: 'vsphere'
        },
      },
      {
        method: :clean_pipeline_jobs,
        params: {
          pipeline_name: 'internetless',
          iaas_type: 'vsphere'
        },
      },
      {
        method: :upgrade_pipeline_jobs,
        params: {
          pipeline_name: 'aws-upgrade',
          iaas_type: 'aws'
        },
      },
      {
        method: :upgrade_pipeline_jobs,
        params: {
          pipeline_name: 'openstack-upgrade',
          iaas_type: 'openstack'
        },
      },
      {
        method: :upgrade_pipeline_jobs,
        params: {
          pipeline_name: 'vsphere-upgrade',
          iaas_type: 'vsphere'
        },
      },
      {
        method: :upgrade_pipeline_jobs,
        params: {
          pipeline_name: 'vcloud-upgrade',
          iaas_type: 'vcloud'
        },
      }
    ].freeze

    def full_suite_pipeline
      full_pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert.yml')))

      PIPELINES.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        full_pipeline_yaml['jobs'].concat(jobs)
      end

      yaml = YAML.dump(full_pipeline_yaml)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6.yml'), yaml)
    end
  end
end
