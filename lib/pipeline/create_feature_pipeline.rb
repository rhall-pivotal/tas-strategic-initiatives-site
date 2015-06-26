require 'mustache'

module Pipeline
  class CreateFeaturePipeline < Mustache
    def initialize(branch_name:, iaas_type:)
      @branch_name = branch_name
      @iaas_type = iaas_type
    end

    def create_pipeline
      dir_name = File.join('ci', 'pipelines', @branch_name)
      FileUtils.mkdir_p(dir_name)

      template_content = File.read(File.join('ci', 'pipelines', 'feature-pipeline-template.yml'))
      File.open(File.join(dir_name, 'pipeline.yml'), 'w') do |f|
        f.write(render(template_content))
      end
    end

    attr_reader :branch_name, :iaas_type
  end
end
