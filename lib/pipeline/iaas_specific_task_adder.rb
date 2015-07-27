require 'yaml'

module Pipeline
  module IaasSpecificTaskAdder
    def add_aws_configure_tasks(pipeline_yaml, template_file)
      extra_config = YAML.load(File.read(File.join(template_directory, template_file)))

      job = pipeline_yaml['jobs'].find { |j| j['name'] =~ /configure-ert/ }

      extra_config.each do |task|
        job['plan'] << task
      end
    end

    def add_vcloud_delete_installation_tasks(pipeline_yaml)
      extra_config = YAML.load(File.read(File.join(template_directory, 'vcloud-delete-installation.yml')))

      destroy_jobs = pipeline_yaml['jobs'].find_all { |j| j['name'] =~ /destroy-environment/ }

      destroy_jobs.each do |job|
        index = job['plan'].find_index { |p| p['task'] == 'destroy' }

        extra_config.each_with_index do |config, i|
          job['plan'].insert(index + i, config)
        end
      end
    end

    def template_directory
      File.join('ci', 'pipelines', 'release', 'template')
    end
  end
end
