require 'mustache'
require 'yaml'

module Pipeline
  class FeaturePipelineCreator < Mustache
    def initialize(branch_name:, iaas_type:)
      @branch_name = branch_name
      @iaas_type = iaas_type
      @ert_version = product_version
      @om_version = @ert_version
    end

    def create_pipeline
      dir_name = make_pipeline_directory

      template_content = File.read(File.join('ci', 'pipelines', 'feature-pipeline-template.yml'))
      write_pipeline_config(dir_name, template_content)
    end

    def create_upgrade_pipeline(ert_initial_full_version:, om_initial_full_version:)
      @ert_initial_full_version = ert_initial_full_version
      @om_initial_full_version = om_initial_full_version

      dir_name = make_pipeline_directory

      template_content = File.read(File.join('ci', 'pipelines', 'feature-upgrade-template.yml'))
      write_pipeline_config(dir_name, template_content)
    end

    def product_version
      handcraft = YAML.load(File.read(File.join('metadata_parts', 'handcraft.yml')))
      fail 'unknown product' unless handcraft['provides_product_versions'][0]['name'] == 'cf'
      full_version = handcraft['provides_product_versions'][0]['version']
      extract_major_minor_version(full_version)
    end

    def ert_initial_version
      extract_major_minor_version(ert_initial_full_version)
    end

    def om_initial_version
      extract_major_minor_version(om_initial_full_version)
    end

    attr_reader :branch_name, :iaas_type, :om_version, :ert_version, :ert_initial_full_version, :om_initial_full_version

    private

    def make_pipeline_directory
      dir_name = File.join('ci', 'pipelines', @branch_name)
      FileUtils.mkdir_p(dir_name)
      dir_name
    end

    def write_pipeline_config(dir_name, template_content)
      File.open(File.join(dir_name, 'pipeline.yml'), 'w') do |f|
        f.write(render(template_content))
      end
    end

    def extract_major_minor_version(full_version)
      full_version.split('.')[0..1].join('.')
    end
  end
end
