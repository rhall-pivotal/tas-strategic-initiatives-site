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
    allow(File).to receive(:read).and_return('')
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
        .and_return('{{branch_name}} {{iaas_type}}')
    )
    expect(file).to receive(:write).with('features/branch aws')

    pipeline_creator.create_pipeline
  end
end
