module Pipeline
  class CreateFeaturePipeline
    def initialize(branch_name:, iaas_type:)
      @branch_name = branch_name
      @iaas_type = iaas_type
    end

    def create_pipeline
      dir_name = File.join('ci', 'pipelines', @branch_name)
      FileUtils.mkdir_p(dir_name)
      File.open(File.join(dir_name, 'pipeline.yml'), 'w') do |f|
        f.write('hello')
      end
    end
  end
end
