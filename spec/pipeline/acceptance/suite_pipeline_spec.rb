require 'spec_helper'
require 'pipeline/suite_pipeline_creator'

describe Pipeline::SuitePipelineCreator do
  subject(:pipeline_creator) do
    Pipeline::SuitePipelineCreator.new
  end

  let(:acceptance_fixture) do
    fixture_path + '/acceptance'
  end

  before do
    allow(File).to(receive(:read).and_call_original)
  end

  context 'when generating a full suite' do
    describe '#full_suite_pipeline' do
      it 'makes the full suite' do
        full_pipeline_fixture = File.join(acceptance_fixture, 'full-pipeline.yml')
        allow(File).to receive(:read).with(full_pipeline_fixture).and_call_original

        expect(File).to receive(:write) do |filename, contents|
          expect(filename).to eq('ci/pipelines/release/ert-1.6.yml')
          expect(YAML.load(contents)).to eq(YAML.load(File.read(full_pipeline_fixture)))
        end

        pipeline_creator.full_suite_pipeline
      end
    end
  end

  context 'when generating a half suite' do
    describe '#half_suite_pipeline' do
      it 'makes the half suite' do
        half_pipeline_fixture = File.join(acceptance_fixture, 'half-pipeline.yml')
        allow(File).to receive(:read).with(half_pipeline_fixture).and_call_original

        expect(File).to receive(:write) do |filename, contents|
          expect(filename).to eq('ci/pipelines/release/ert-1.6-half.yml')
          expect(YAML.load(contents)).to eq(YAML.load(File.read(half_pipeline_fixture)))
        end

        pipeline_creator.half_suite_pipeline
      end
    end
  end
end
