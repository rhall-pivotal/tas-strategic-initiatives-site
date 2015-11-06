require 'yaml'

module Pipeline
  module IaasSpecificTaskAdder
    def fetch_configure_tasks(task_key, *template_files)
      tasks = { task_key => [] }

      template_files.each do |template|
        tasks[task_key] << { task: File.read(File.join(template_directory, template)) }
      end

      tasks
    end

    def fetch_verify_internetless_job
      fetch_configure_tasks(:verify_internetless_plan, 'internetless-verification.yml')
    end

    def template_directory
      File.join('ci', 'pipelines', 'release', 'template')
    end
  end
end
