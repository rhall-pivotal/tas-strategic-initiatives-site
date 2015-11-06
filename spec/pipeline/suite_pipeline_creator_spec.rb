require 'spec_helper'
require 'pipeline/suite_pipeline_creator'

describe Pipeline::SuitePipelineCreator do
  subject(:pipeline_creator) do
    Pipeline::SuitePipelineCreator.new
  end

  let(:template) do
    <<YAML
---
jobs:
  - name: destroy-environment-{{pipeline_name}}
    plan:
      - get: environment
        resource: environment-{{environment_pool}}
        trigger: true
        passed: [claim-environment-{{pipeline_name}}]
      - task: destroy
        tags: [{{iaas_type}}]
        file: p-runtime/ci/jobs/destroy-environment.yml
  - name: deploy-ops-manager-{{pipeline_name}}
    plan:
      - task: some-task
  {{#verify_internetless_plan}}
  {{task}}
  {{/verify_internetless_plan}}
  - name: configure-microbosh-{{pipeline_name}}
    plan:
      - get: p-runtime
        passed:
        {{#verify_internetless_plan}}
        - verify-internetless
        {{/verify_internetless_plan}}
        {{^verify_internetless_plan}}
        - deploy-ops-manager-{{pipeline_name}}
        {{/verify_internetless_plan}}
      - get: environment
        resource: environment-{{pipeline_name}}
        passed:
        {{#verify_internetless_plan}}
        - verify-internetless
        {{/verify_internetless_plan}}
        {{^verify_internetless_plan}}
        - deploy-ops-manager-{{pipeline_name}}
        {{/verify_internetless_plan}}
  - name: configure-ert-{{pipeline_name}}
    plan:
      - task: some-task
        tags: [{{iaas_type}}]
        resource: environment-{{environment_pool}}
      {{#aws_configure_tasks}}
      {{task}}
      {{/aws_configure_tasks}}
  - name: another-job-{{pipeline_name}}
    plan:
      - task: another-task
        tags: [{{iaas_type}}]
        resource: environment-{{environment_pool}}
  - name: destroy-environment-final-{{pipeline_name}}
    plan:
      - get: environment
        resource: environment-{{environment_pool}}
        trigger: true
        passed: [claim-environment-{{pipeline_name}}]
      - task: destroy
        tags: [{{iaas_type}}]
        file: p-runtime/ci/jobs/destroy-environment.yml
  - name: release-environment-{{pipeline_name}}
    plan:
      - get: environment
        resource: environment-{{environment_pool}}
        trigger: true
        passed: [destroy-environment-{{pipeline_name}}]
      - task: release
        tags: [{{iaas_type}}]
        file: p-runtime/ci/jobs/release-environment.yml
YAML
  end

  let(:aws_extra_config) do
    <<YAML
- task: some-aws-task
        tags: aws
YAML
  end

  let(:aws_experimental_config) do
    <<YAML
- task: some-experimental-aws-task
        tags: aws
YAML
  end

  let(:aws_extra_config_upgrade) do
    <<YAML
- task: some-aws-upgrade-task
        tags: aws
YAML
  end

  let(:verify_internetless_config) do
    <<YAML
- name: verify-internetless
    plan:
    - task: verify
YAML
  end

  let(:multi_az) do
    <<YAML
- name: multi-az-specific-job
    plan:
    - task: do-multi-az-stuff
YAML
  end

  before do
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/clean.yml')
        .and_return(template)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/upgrade.yml')
        .and_return(template)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/aws-external-config.yml')
        .and_return(aws_extra_config)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/aws-enable-experimental.yml')
        .and_return(aws_experimental_config)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/aws-external-config-upgrade.yml')
        .and_return(aws_extra_config_upgrade)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/internetless-verification.yml')
        .and_return(verify_internetless_config)
    )
    allow(File).to(
      receive(:read)
        .with('ci/pipelines/release/template/multi-az.yml')
        .and_return(multi_az)
    )
  end

  it 'has a constructor that takes no arguments' do
    expect(pipeline_creator).to be_a(Pipeline::SuitePipelineCreator)
  end

  describe '#iaas_specific_pipeline_jobs' do
    context 'when the iaas is vsphere' do
      context 'when the pipeline_name is something other than internetless' do
        it 'returns empty hash' do
          additional_jobs = pipeline_creator.iaas_specific_pipeline_job(
            pipeline_type: :upgrade, params: { pipeline_name: 'vsphere-clean', iaas_type: 'vsphere' }
          )

          expect(additional_jobs).to be_nil
        end
      end

      context 'when the pipeline_name is internetless' do
        it 'returns the internetless job' do
          additional_jobs = pipeline_creator.iaas_specific_pipeline_job(
            pipeline_type: :clean, params: { pipeline_name: 'internetless', iaas_type: 'vsphere' })
          expect(additional_jobs[:verify_internetless_plan][0][:task])
            .to eq(verify_internetless_config)
        end
      end

      context 'when the environment_pool is multi-az' do
        it 'returns the multi-az job' do
          additional_jobs = pipeline_creator.iaas_specific_pipeline_job(
            pipeline_type: :clean, params: {
              pipeline_name: 'vsphere-upgrade',
              iaas_type: 'vsphere',
              environment_pool: 'multi-az'
            })
          expect(additional_jobs[:multi_az][0][:task])
            .to eq(multi_az)
        end
      end
    end

    context 'when the iaas is aws' do
      context 'when the method is clean' do
        it 'returns only the aws-external-config and aws-enable-experimental jobs' do
          additional_jobs = pipeline_creator.iaas_specific_pipeline_job(
            pipeline_type: :clean, params: { pipeline_name: 'aws-clean', iaas_type: 'aws' })

          expect(additional_jobs[:aws_configure_tasks][0][:task])
            .to eq(aws_extra_config)

          expect(additional_jobs[:aws_configure_tasks][1][:task])
            .to eq(aws_experimental_config)
        end
      end

      context 'when the method is upgrade' do
        it 'returns only the aws-external-config-upgrade' do
          additional_jobs = pipeline_creator.iaas_specific_pipeline_job(
            pipeline_type: :upgrade, params: { pipeline_name: 'aws-clean', iaas_type: 'aws' }
          )

          expect(additional_jobs[:aws_configure_tasks][0][:task])
            .to eq(aws_extra_config_upgrade)
        end
      end
    end
  end
end
