require 'mustache'
require 'yaml'
require 'pipeline/iaas_specific_task_adder'

module Pipeline
  class SuitePipelineCreator < Mustache
    include IaasSpecificTaskAdder

    HALF_PIPELINES = [
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'aws-clean', iaas_type: 'aws' } },
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'internetless', iaas_type: 'vsphere' } },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'aws-upgrade', iaas_type: 'aws' } },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'vsphere-upgrade', iaas_type: 'vsphere' } },
    ].freeze

    FULL_PIPELINES = [
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

    def environment_pool
      case pipeline_name
      when 'internetless'
        pipeline_name
      when 'aws-upgrade'
        'aws-east'
      else
        iaas_type
      end
    end

    def upgrade_pipeline_jobs(pipeline_name:, iaas_type:)
      pipeline_yaml = pipeline_jobs(
        pipeline_name: pipeline_name,
        iaas_type: iaas_type,
        template_path: File.join(template_directory, 'upgrade.yml')
      )

      add_aws_configure_tasks(pipeline_yaml, 'aws-external-config-upgrade.yml') if iaas_type == 'aws'

      pipeline_yaml
    end

    def clean_pipeline_jobs(pipeline_name:, iaas_type:)
      pipeline_yaml = pipeline_jobs(
        pipeline_name: pipeline_name,
        iaas_type: iaas_type,
        template_path: File.join(template_directory, 'clean.yml')
      )

      add_aws_configure_tasks(pipeline_yaml, 'aws-external-config.yml') if iaas_type == 'aws'
      add_verify_internetless_job(pipeline_yaml) if pipeline_name == 'internetless'

      pipeline_yaml
    end

    def half_suite_pipeline
      half_pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert-half.yml')))

      HALF_PIPELINES.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        half_pipeline_yaml['jobs'].concat(jobs)
      end

      yaml = YAML.dump(half_pipeline_yaml)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6-half.yml'), yaml)
    end

    def full_suite_pipeline
      full_pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert.yml')))

      FULL_PIPELINES.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        full_pipeline_yaml['jobs'].concat(jobs)
      end

      yaml = YAML.dump(full_pipeline_yaml)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6.yml'), yaml)
    end

    attr_reader :pipeline_name, :iaas_type

    private

    def pipeline_jobs(pipeline_name:, iaas_type:, template_path:)
      @pipeline_name = pipeline_name
      @iaas_type = iaas_type

      pipeline_yaml = YAML.load(render(File.read(template_path)))

      add_vcloud_delete_installation_tasks(pipeline_yaml) if iaas_type == 'vcloud'
      pipeline_yaml
    end
  end
end
