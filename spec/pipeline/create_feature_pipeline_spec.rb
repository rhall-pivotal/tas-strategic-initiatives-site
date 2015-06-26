require 'spec_helper'
require 'pipeline/create_feature_pipeline'

describe Pipeline::CreateFeaturePipeline do
  subject(:pipeline_creator) do
    Pipeline::CreateFeaturePipeline.new(
      branch_name: 'features/branch',
      iaas_type: 'aws',
    )
  end
  let(:file) { instance_double(File) }
  let(:handcraft) do
    <<YAML
---
metadata_version: '1.5'
provides_product_versions:
- name: cf
  version: 1.5.0.0
product_version: &product_version "1.5.0.0$PRERELEASE_VERSION$"
label: Pivotal Elastic Runtime
YAML
  end

  before do
    allow(File).to receive(:open).and_yield(file)
    allow(file).to receive(:write)
    allow(FileUtils).to receive(:mkdir_p)
    allow(File).to receive(:read).and_return('')
    allow(File).to receive(:read).and_return(handcraft)
  end

  it 'has a constructor that takes two arguments' do
    expect(pipeline_creator).to be_a(Pipeline::CreateFeaturePipeline)
  end

  it 'creates the directory for the pipeline config file' do
    expect(FileUtils).to receive(:mkdir_p).with('ci/pipelines/features/branch')

    pipeline_creator.create_pipeline
  end

  it 'creates the pipeline configuration file with the correct name' do
    expect(File).to receive(:open).with('ci/pipelines/features/branch/pipeline.yml', 'w')

    pipeline_creator.create_pipeline
  end

  it 'puts the rendered content from the template into the configuration file' do
    expect(File).to(
      receive(:read)
        .with('ci/pipelines/feature-pipeline-template.yml')
        .and_return('{{branch_name}} {{iaas_type}} {{om_version}} {{ert_version}}')
    )
    expect(file).to receive(:write).with('features/branch aws 1.5 1.5')

    pipeline_creator.create_pipeline
  end

  it 'returns the product version from the handcraft.yml file' do
    expect(File).to receive(:read).with('metadata_parts/handcraft.yml').and_return(handcraft)
    expect(pipeline_creator.product_version).to eq('1.5')
  end

  context 'when the first product is not cf' do
    let(:handcraft) do
      <<YAML
---
metadata_version: '1.5'
provides_product_versions:
- name: foo
  version: 1.5.0.0
product_version: &product_version "1.5.0.0$PRERELEASE_VERSION$"
label: Pivotal Elastic Runtime
YAML
    end

    it 'raises' do
      expect { pipeline_creator.product_version }.to raise_error
    end
  end
end
