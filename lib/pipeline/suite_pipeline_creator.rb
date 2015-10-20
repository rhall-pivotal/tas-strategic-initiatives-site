require 'mustache'
require 'yaml'
require 'pipeline/iaas_specific_task_adder'

module Pipeline
  class SuitePipelineCreator < Mustache
    include IaasSpecificTaskAdder

    HALF_PIPELINES = [
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'aws-clean', iaas_type: 'aws' } },
      {
        method: :clean_pipeline_jobs,
        params: { pipeline_name: 'internetless', iaas_type: 'vsphere' },
        group_name: 'vsphere-internetless'
      },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'aws-upgrade', iaas_type: 'aws' } },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'vsphere-upgrade', iaas_type: 'vsphere' } },
    ].freeze

    FULL_PIPELINES = [
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'aws-clean', iaas_type: 'aws' } },
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'openstack-clean', iaas_type: 'openstack' } },
      { method: :clean_pipeline_jobs, params: { pipeline_name: 'vsphere-clean', iaas_type: 'vsphere' } },
      { method: :clean_pipeline_jobs,
        params: { pipeline_name: 'internetless', iaas_type: 'vsphere' },
        group_name: 'vsphere-internetless'
      },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'aws-upgrade', iaas_type: 'aws' } },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'vsphere-upgrade', iaas_type: 'vsphere' } },
      { method: :upgrade_pipeline_jobs, params: { pipeline_name: 'vcloud-upgrade', iaas_type: 'vcloud' } }
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
      add_aws_configure_tasks(pipeline_yaml, 'aws-enable-experimental.yml') if iaas_type == 'aws'
      add_verify_internetless_job(pipeline_yaml) if pipeline_name == 'internetless'

      pipeline_yaml
    end

    def half_suite_pipeline
      yaml = create_pipeline_yaml(HALF_PIPELINES)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6-half.yml'), yaml)
    end

    def full_suite_pipeline
      yaml = create_pipeline_yaml(FULL_PIPELINES)

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6.yml'), yaml)
    end

    attr_reader :pipeline_name, :iaas_type

    private

    def critical_jobs(pipeline_yaml)
      pipeline_yaml['jobs']
        .select { |j| j['name'].start_with?('release-environment') }
        .reject { |j| j['name'].include?('openstack') }
        .map { |j| j['name'] }
    end

    def step_needing_passed_criteria(pipeline_yaml)
      pipeline_yaml['jobs'].find { |j| j['name'] == 'promote-ert' }['plan'].find { |s| s['get'] == 'p-runtime' }
    end

    def pipeline_jobs(pipeline_name:, iaas_type:, template_path:)
      @pipeline_name = pipeline_name
      @iaas_type = iaas_type

      pipeline_yaml = YAML.load(render(File.read(template_path)))

      add_vcloud_delete_installation_tasks(pipeline_yaml) if iaas_type == 'vcloud'
      pipeline_yaml
    end

    def create_pipeline_yaml(pipelines)
      pipeline_yaml = YAML.load(File.read(File.join(template_directory, 'ert.yml')))

      pipelines.each do |config|
        jobs = send(config[:method], config[:params])['jobs']
        pipeline_yaml['jobs'].concat(jobs)
      end

      step_needing_passed_criteria(pipeline_yaml)['passed'] = critical_jobs(pipeline_yaml)

      groups = pipeline_groups(pipeline_yaml, pipelines)

      pipeline_yaml['groups'] = groups

      YAML.dump(pipeline_yaml)
    end

    def pipeline_groups(pipeline_yaml, pipelines)
      groups = ['common'] + pipelines.map do |p|
        p[:group_name] || p[:params][:pipeline_name]
      end

      pipeline_yaml['groups'].select { |g| groups.include?(g['name']) }
    end
  end
end
