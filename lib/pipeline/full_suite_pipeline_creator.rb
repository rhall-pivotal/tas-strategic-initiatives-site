require 'mustache'
require 'yaml'

module Pipeline
  class FullSuitePipelineCreator < Mustache
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
      pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert.yml')))

      PIPELINES.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        pipeline_yaml['jobs'].concat(jobs)
      end

      yaml = YAML.dump(pipeline_yaml)
      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.5.yml'), yaml)
    end

    def clean_pipeline_jobs(pipeline_name:, iaas_type:)
      pipeline_yaml = pipeline_jobs(
        pipeline_name: pipeline_name,
        iaas_type: iaas_type,
        template_path: File.join(template_directory,  'clean.yml')
      )

      add_aws_configure_tasks(pipeline_yaml, 'aws-external-config.yml') if iaas_type == 'aws'

      pipeline_yaml
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

    attr_reader :pipeline_name, :iaas_type

    private

    def pipeline_jobs(pipeline_name:, iaas_type:, template_path:)
      @pipeline_name = pipeline_name
      @iaas_type = iaas_type

      pipeline_yaml = YAML.load(render(File.read(template_path)))

      add_vcloud_delete_installation_tasks(pipeline_yaml) if iaas_type == 'vcloud'
      pipeline_yaml
    end

    def add_aws_configure_tasks(pipeline_yaml, template_file)
      extra_config = YAML.load(File.read(File.join(template_directory, template_file)))

      job = pipeline_yaml['jobs'].find { |j| j['name'] =~ /configure-ert/ }

      extra_config.each do |task|
        job['plan'] << task
      end
    end

    def add_vcloud_delete_installation_tasks(pipeline_yaml)
      extra_config = YAML.load(File.read(File.join(template_directory, 'vcloud-delete-installation.yml')))

      job = pipeline_yaml['jobs'].find { |j| j['name'] =~ /destroy-environment/ }
      index = job['plan'].find_index { |p| p['task'] == 'destroy' }

      extra_config.each_with_index do |config, i|
        job['plan'].insert(index + i, config)
      end
    end

    def template_directory
      File.join('ci', 'pipelines', 'release', 'template')
    end
  end
end
