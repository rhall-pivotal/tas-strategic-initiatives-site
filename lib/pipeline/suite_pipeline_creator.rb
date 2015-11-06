require 'mustache'
require 'yaml'
require 'pipeline/iaas_specific_task_adder'

module Pipeline
  class SuitePipelineCreator < Mustache
    include IaasSpecificTaskAdder

    HALF_PIPELINES = [
      { pipeline_type: :clean,
        params: { pipeline_name: 'aws-clean', iaas_type: 'aws', environment_pool: 'aws' }
      },
      {
        pipeline_type: :clean,
        params: { pipeline_name: 'internetless', iaas_type: 'vsphere', environment_pool: 'internetless' },
        group_name: 'vsphere-internetless'
      },
      { pipeline_type: :upgrade,
        params: { pipeline_name: 'aws-upgrade', iaas_type: 'aws', environment_pool: 'aws' }
      },
      { pipeline_type: :upgrade,
        params: { pipeline_name: 'vsphere-upgrade', iaas_type: 'vsphere', environment_pool: 'vsphere' }
      },
    ].freeze

    MULTI_AZ_PIPELINES = [
      { pipeline_type: :clean,
        params: { pipeline_name: 'vsphere-clean', iaas_type: 'vsphere', environment_pool: 'multi-az' }
      },
    ].freeze

    FULL_PIPELINES = [
      {
        pipeline_type: :clean,
        params: { pipeline_name: 'aws-clean', iaas_type: 'aws', environment_pool: 'aws' }
      },
      {
        pipeline_type: :clean,
        params: { pipeline_name: 'openstack-clean', iaas_type: 'openstack', environment_pool: 'openstack' }
      },
      { pipeline_type: :clean,
        params: { pipeline_name: 'vsphere-clean', iaas_type: 'vsphere', environment_pool: 'vsphere' }
      },
      { pipeline_type: :clean,
        params: { pipeline_name: 'internetless', iaas_type: 'vsphere', environment_pool: 'internetless' },
        group_name: 'vsphere-internetless'
      },
      { pipeline_type: :upgrade,
        params: { pipeline_name: 'aws-upgrade', iaas_type: 'aws', environment_pool: 'aws' }
      },
      { pipeline_type: :upgrade,
        params: { pipeline_name: 'vsphere-upgrade', iaas_type: 'vsphere', environment_pool: 'vsphere' }
      },
      { pipeline_type: :upgrade,
        params: { pipeline_name: 'vcloud-upgrade', iaas_type: 'vcloud', environment_pool: 'vcloud' }
      }
    ].freeze

    def iaas_specific_pipeline_job(pipeline_type:, params:, group_name: nil)
      _group_name = group_name
      if params[:iaas_type] == 'aws'
        aws_stuff(pipeline_type)
      elsif params[:pipeline_name] == 'internetless'
        fetch_verify_internetless_job
      elsif params[:iaas_type] == 'vsphere' && params[:environment_pool] == 'multi-az'
        multi_az_stuff(pipeline_type)
      end
    end

    def half_suite_pipeline
      yaml = create_pipeline_yaml(HALF_PIPELINES, 'ert.yml')

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6-half.yml'), yaml)
    end

    def full_suite_pipeline
      yaml = create_pipeline_yaml(FULL_PIPELINES, 'ert.yml')

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6.yml'), yaml)
    end

    def multi_az_pipeline
      yaml = create_pipeline_yaml(MULTI_AZ_PIPELINES, 'ert-multi-az.yml')

      File.write(File.join('ci', 'pipelines', 'release', 'ert-1.6-multi-az.yml'), yaml)
    end

    attr_reader :pipeline_name, :iaas_type, :environment_pool

    private

    def multi_az_stuff(pipeline_type)
      case pipeline_type
      when :upgrade
        fail 'multi-az not supported in 1.5 - fix me for 1.6->1.7 upgrade'
      when :clean
        fetch_configure_tasks(:multi_az, 'multi-az.yml')
      end
    end

    def aws_stuff(pipeline_type)
      case pipeline_type
      when :upgrade
        fetch_configure_tasks(:aws_configure_tasks, 'aws-external-config-upgrade.yml')
      when :clean
        fetch_configure_tasks(:aws_configure_tasks, 'aws-external-config.yml', 'aws-enable-experimental.yml')
      end
    end

    def critical_jobs(pipeline_yaml)
      pipeline_yaml['jobs']
        .select { |j| j['name'].start_with?('release-environment') }
        .reject { |j| j['name'].include?('openstack') }
        .map { |j| j['name'] }
    end

    def step_needing_passed_criteria(pipeline_yaml)
      pipeline_yaml['jobs'].find { |j| j['name'] == 'promote-ert' }['plan'].find { |s| s['get'] == 'p-runtime' }
    end

    def construct_template_path(method)
      file = method == :upgrade ? 'upgrade.yml' : 'clean.yml'

      File.join(template_directory, file)
    end

    def pipeline_jobs(pipeline_type, additional_jobs, pipeline_name:, iaas_type:, environment_pool:)
      @pipeline_name = pipeline_name
      @iaas_type = iaas_type
      @environment_pool = environment_pool

      template = File.read(construct_template_path(pipeline_type))

      YAML.load(render(template, additional_jobs))
    end

    def create_pipeline_yaml(pipelines, main_template_file)
      pipeline_yaml = YAML.load(File.read(File.join(template_directory, main_template_file)))

      pipelines.each do |config|
        params = config[:params]

        iaas_specific_jobs = iaas_specific_pipeline_job(config)
        iaas_specific_yaml = pipeline_jobs(config[:pipeline_type], iaas_specific_jobs, params)

        pipeline_yaml['jobs'].concat(iaas_specific_yaml['jobs'])
      end

      if main_template_file == 'ert.yml'
        step_needing_passed_criteria(pipeline_yaml)['passed'] = critical_jobs(pipeline_yaml)
      end

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
