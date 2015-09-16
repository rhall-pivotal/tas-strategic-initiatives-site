require 'spec_helper'
require 'pipeline/feature_pipeline_creator'
require 'pipeline/suite_pipeline_creator'

describe 'Generated pipelines' do
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
    allow(FileUtils).to receive(:mkdir_p)
  end

  context 'full suite pipeline' do
    context 'aws clean' do
      it 'contains configure experimental features task' do
        expect(File).to(
          receive(:write) do |filename, content|
            expect(filename).to eq('ci/pipelines/release/ert-1.6-half.yml')
            pipeline = YAML.load(content)
            jobs = pipeline['jobs']
            configure_ert_job = jobs.find { |j| j['name'] == 'configure-ert-aws-clean' }
            configure_tasks = configure_ert_job['plan'].select { |s| !s['task'].nil? }
            configure_task_names = configure_tasks.map { |t| t['task'] }
            expect(configure_task_names).to include('configure-experimental-features')
          end
        )

        Pipeline::SuitePipelineCreator.new.half_suite_pipeline
      end
    end

    context 'aws upgrade' do
      it 'does not contain configure experimental features task' do
        expect(File).to(
          receive(:write) do |filename, content|
            expect(filename).to eq('ci/pipelines/release/ert-1.6-half.yml')
            pipeline = YAML.load(content)
            jobs = pipeline['jobs']
            configure_ert_job = jobs.find { |j| j['name'] == 'configure-ert-aws-upgrade' }
            configure_tasks = configure_ert_job['plan'].select { |s| !s['task'].nil? }
            configure_task_names = configure_tasks.map { |t| t['task'] }
            expect(configure_task_names).not_to include('configure-experimental-features')
          end
        )

        Pipeline::SuitePipelineCreator.new.half_suite_pipeline
      end
    end
  end

  context 'half suite pipeline' do
    context 'aws clean' do
      it 'contains configure experimental features task' do
        expect(File).to(
          receive(:write) do |filename, content|
            expect(filename).to eq('ci/pipelines/release/ert-1.6-half.yml')
            pipeline = YAML.load(content)
            jobs = pipeline['jobs']
            configure_ert_job = jobs.find { |j| j['name'] == 'configure-ert-aws-clean' }
            configure_tasks = configure_ert_job['plan'].select { |s| !s['task'].nil? }
            configure_task_names = configure_tasks.map { |t| t['task'] }
            expect(configure_task_names).to include('configure-experimental-features')
          end
        )

        Pipeline::SuitePipelineCreator.new.half_suite_pipeline
      end
    end
  end
  context 'clean-install feature pipeline' do
    context 'on aws' do
      it 'does not configure experimental features' do
        expect(File).to receive(:open).with('ci/pipelines/branch/pipeline.yml', 'w').and_yield(file)

        expect(file).to(
          receive(:write) do |content|
            pipeline = YAML.load(content)
            jobs = pipeline['jobs']
            configure_ert_job = jobs.find { |j| j['name'] == 'configure-ert' }
            configure_tasks = configure_ert_job['plan'].select { |s| !s['task'].nil? }
            configure_task_names = configure_tasks.map { |t| t['task'] }
            expect(configure_task_names).not_to include('configure-experimental-features')
          end
        )
        Pipeline::FeaturePipelineCreator.new(
          branch_name: 'branch',
          iaas_type: 'aws',
        ).create_pipeline
      end
    end
  end
  context 'upgrade feature pipeline'
end
