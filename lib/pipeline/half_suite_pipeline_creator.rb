require 'pipeline/suite_pipeline_creator.rb'

module Pipeline
  class HalfSuitePipelineCreator < Pipeline::SuitePipelineCreator
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
          pipeline_name: 'vsphere-upgrade',
          iaas_type: 'vsphere'
        },
      },
    ].freeze

    def half_suite_pipeline
      half_pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert-half.yml')))

      PIPELINES.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        half_pipeline_yaml['jobs'].concat(jobs)
      end

      yaml = YAML.dump(half_pipeline_yaml)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6-half.yml'), yaml)
    end
  end
end
