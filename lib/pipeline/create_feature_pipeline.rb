require 'mustache'
require 'yaml'

module Pipeline
  class CreateFeaturePipeline < Mustache
    def initialize(branch_name:, iaas_type:)
      @branch_name = branch_name
      @iaas_type = iaas_type
      @ert_version = product_version
      @om_version = @ert_version
    end

    def create_pipeline
      dir_name = File.join('ci', 'pipelines', @branch_name)
      FileUtils.mkdir_p(dir_name)

      template_content = File.read(File.join('ci', 'pipelines', 'feature-pipeline-template.yml'))
      File.open(File.join(dir_name, 'pipeline.yml'), 'w') do |f|
        f.write(render(template_content))
      end
    end

    def product_version
      handcraft = YAML.load(File.read(File.join('metadata_parts', 'handcraft.yml')))
      fail 'unknown product' unless handcraft['provides_product_versions'][0]['name'] == 'cf'
      full_version = handcraft['provides_product_versions'][0]['version']
      full_version.split('.')[0..1].join('.')
    end

    attr_reader :branch_name, :iaas_type, :om_version, :ert_version
  end
end
