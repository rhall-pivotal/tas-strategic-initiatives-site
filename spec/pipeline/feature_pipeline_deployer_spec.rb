require 'spec_helper'
require 'pipeline/feature_pipeline_deployer'

describe Pipeline::FeaturePipelineDeployer do
  subject(:pipeline_deployer) do
    Pipeline::FeaturePipelineDeployer.new(
      branch_name: branch_name
    )
  end
  let(:branch_name) { 'features/some-branch-name' }
  let(:colon_delimited_branch_name) { 'features::some-branch-name' }
  let(:config_file_path) { "ci/pipelines/#{branch_name}/pipeline.yml" }
  let(:fly_command) { "fly -t ci configure #{branch_name} -c ci/pipelines/feature-branch.yml" }

  before do
    allow(File).to receive(:join).with('ci', 'pipelines', branch_name, 'pipeline.yml')
      .and_return(config_file_path)
    allow(File).to receive(:exist?).with(config_file_path).and_return(true)
    allow(pipeline_deployer).to receive(:system)
  end

  it 'has a constructor that takes one argument' do
    expect(pipeline_deployer).to be_a(Pipeline::FeaturePipelineDeployer)
  end

  it 'uses the pipeline configuration file with the same name as the branch' do
    expect(File).to receive(:exist?).with(config_file_path).and_return(true)

    pipeline_deployer.deploy_pipeline
  end

  context 'when the configuration file does not exist' do
    let(:branch_name) { 'some-fake-branch-name' }

    before do
      expect(File).to receive(:exist?).with(config_file_path).and_return(false)
    end

    it 'raises an error' do
      expect { pipeline_deployer.deploy_pipeline }
        .to raise_error(
          Pipeline::FeaturePipelineDeployer::NoConfigFileError,
          "Unable to find pipeline configuration for #{branch_name}"
        )
    end
  end

  context 'shelling out to fly' do
    it 'converts slashes in the branch name to double colons' do
      expect(pipeline_deployer).to receive(:system)
        .with(/fly -t ci configure #{colon_delimited_branch_name}/)

      pipeline_deployer.deploy_pipeline
    end

    it 'uses the configuration file with the same branch name' do
      expect(pipeline_deployer).to receive(:system)
        .with("fly -t ci configure #{colon_delimited_branch_name} -c ci/pipelines/#{branch_name}/pipeline.yml")

      pipeline_deployer.deploy_pipeline
    end
  end
end
