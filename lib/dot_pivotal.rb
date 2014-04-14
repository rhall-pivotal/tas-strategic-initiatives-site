require 'vara/product_metadata_builder'
require 'vara/product_metadata'

class DotPivotal
  attr_reader :base_dir, :build_identifier, :metadata_file

  def initialize(base_dir, build_identifier = 'local')
    @base_dir = base_dir
    @build_identifier = build_identifier
  end

  def build_metadata
    @metadata_file = Vara::ProductMetadataBuilder.build(base_dir)
  end

  def stemcell_version
    @stemcell_version ||= product_metadata.stemcell_metadata.version
  end

  def metadata_ref
    `cd #{base_dir} && git log -1 --format=format:%h`
  end

  def releases
    product_metadata.releases_metadata.map do |release_metadata|
      {
        'name' => release_metadata['name'],
        'version' => release_metadata['version']
      }
    end
  end

  private

  def product_metadata
    @product_metadata ||= Vara::ProductMetadata.from_file(metadata_file)
  end
end
