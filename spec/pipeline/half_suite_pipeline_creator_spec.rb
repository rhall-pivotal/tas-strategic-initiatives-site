require 'spec_helper'
require 'pipeline/half_suite_pipeline_creator'

describe Pipeline::HalfSuitePipelineCreator do
  subject(:pipeline_creator) do
    Pipeline::HalfSuitePipelineCreator.new
  end
  let(:ert_general) do
    <<YAML
---
resources:
  - name: some-resource
    type: git
  - name: another-resource
    type: s3
    source:
      some_key: just-a-key
jobs:
  - name: a-generic-job
    plan:
      - get: a-get-task
        resource: some-resource
      - task: a-generic-task
YAML
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
  - name: configure-microbosh-{{pipeline_name}}
    plan:
      - get: p-runtime
        passed: [deploy-ops-manager-{{pipeline_name}}]
      - get: environment
        resource: environment-{{pipeline_name}}
        passed: [deploy-ops-manager-{{pipeline_name}}]
  - name: configure-ert-{{pipeline_name}}
    plan:
      - task: some-task
        tags: [{{iaas_type}}]
        resource: environment-{{environment_pool}}
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
YAML
  end

  let(:aws_extra_config) do
    <<YAML
- task: some-aws-task
  tags: aws
YAML
  end

  let(:aws_extra_config_upgrade) do
    <<YAML
- task: some-aws-upgrade-task
  tags: aws
YAML
  end

  let(:vcloud_extra_config) do
    <<YAML
- task: some-vcloud-task
  tags: vcloud
- task: another-vcloud-task
  tags: vcloud
YAML
  end

  let(:verify_internetless_config) do
    <<YAML
---
name: verify-internetless
plan:
- task: verify
YAML
  end

  before do
    allow(File).to receive(:read).and_return(template)
    allow(File).to receive(:read).with('ci/pipelines/release/template/ert-half.yml')
      .and_return(ert_general)
    allow(File).to receive(:read).with('ci/pipelines/release/template/aws-external-config.yml')
      .and_return(aws_extra_config)
    allow(File).to receive(:read).with('ci/pipelines/release/template/aws-external-config-upgrade.yml')
      .and_return(aws_extra_config_upgrade)
    allow(File).to receive(:read).with('ci/pipelines/release/template/vcloud-delete-installation.yml')
      .and_return(vcloud_extra_config)
    allow(File).to receive(:read).with('ci/pipelines/release/template/internetless-verification.yml')
      .and_return(verify_internetless_config)
  end

  describe '#half_suite_pipeline' do
    it 'makes the half suite' do
      half_pipeline_fixture = File.join(fixture_path, 'half-pipeline.yml')
      allow(File).to receive(:read).with(half_pipeline_fixture).and_call_original

      expect(File).to receive(:write) do |filename, contents|
        expect(filename).to eq('ci/pipelines/release/ert-1.6-half.yml')
        expect(contents).to eq(File.read(half_pipeline_fixture))
      end

      pipeline_creator.half_suite_pipeline
    end
  end
end
