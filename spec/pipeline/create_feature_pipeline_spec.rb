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

  before do
    allow(File).to receive(:open).and_yield(file)
    allow(file).to receive(:write)
    allow(FileUtils).to receive(:mkdir_p)
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

  it 'puts some content into the configuration file' do
    expect(file).to receive(:write).with('hello')

    pipeline_creator.create_pipeline
  end
end
