module Pipeline
  class FeaturePipelineDeployer
    class NoConfigFileError < StandardError; end

    def initialize(branch_name:)
      @branch_name = branch_name
    end

    def deploy_pipeline
      file_path = File.join('ci', 'pipelines', branch_name, 'pipeline.yml')

      unless File.exist?(file_path)
        fail NoConfigFileError, "Unable to find pipeline configuration for #{branch_name}"
      end

      colon_branch_name = branch_name.sub('/', '::')
      system("fly -t ci configure #{colon_branch_name} -c ci/pipelines/#{branch_name}/pipeline.yml")
    end

    attr_reader :branch_name
  end
end
