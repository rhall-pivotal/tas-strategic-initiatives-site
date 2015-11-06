require 'spec_helper'
require 'yaml'
require 'pipeline/iaas_specific_task_adder'

describe Pipeline::IaasSpecificTaskAdder do
  subject(:task_adder) do
    Class.new { include Pipeline::IaasSpecificTaskAdder }.new
  end

  let(:aws_extra_config) do
    <<YAML
- task: some-aws-task
  tags: aws
YAML
  end

  let(:another_aws_config) do
    <<YAML
- task: another-aws-task
  tags: aws
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

  describe '#fetch_aws_configure_tasks' do
    before do
      allow(File).to(
        receive(:read).with('ci/pipelines/release/template/aws-extra-config.yml')
          .and_return(aws_extra_config)
      )

      allow(File).to(
        receive(:read).with('ci/pipelines/release/template/another-aws-config.yml')
          .and_return(another_aws_config)
      )
    end

    it 'adds all tasks extra template file to configure-ert job' do
      aws_tasks = task_adder.fetch_configure_tasks(:aws_configure_tasks, 'aws-extra-config.yml', 'another-aws-config.yml')

      expect(aws_tasks[:aws_configure_tasks][0][:task]).to eq(aws_extra_config)
      expect(aws_tasks[:aws_configure_tasks][1][:task]).to eq(another_aws_config)
    end
  end

  describe '#fetch_verify_internetless_job' do
    before do
      allow(File).to(
        receive(:read).with('ci/pipelines/release/template/internetless-verification.yml')
          .and_return(verify_internetless_config)
      )
    end

    it 'adds the verify_internetless job after deploy-ops-manager' do
      internetless_plan = task_adder.fetch_verify_internetless_job
      expect(internetless_plan[:verify_internetless_plan][0][:task]).to eq(verify_internetless_config)
    end
  end
end
