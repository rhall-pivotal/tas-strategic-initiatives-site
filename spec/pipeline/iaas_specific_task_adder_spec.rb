require 'spec_helper'
require 'yaml'
require 'pipeline/iaas_specific_task_adder'

describe Pipeline::IaasSpecificTaskAdder do
  subject(:task_adder) do
    Class.new { include Pipeline::IaasSpecificTaskAdder }.new
  end

  let(:pipeline_yaml) do
    YAML.load(<<YAML)
---
jobs:
  - name: destroy-environment
    plan:
      - get: p-runtime
        passed: [claim-environment]
      - get: environment
        passed: [claim-environment]
      - task: destroy
        tags: [some-iaas]
  - name: deploy-ops-manager
    plan:
      - task: some-task
  - name: configure-microbosh
    plan:
      - get: p-runtime
        passed: [deploy-ops-manager]
      - get: environment
        passed: [deploy-ops-manager]
  - name: configure-ert
    plan:
      - task: some-task
        tags: [some-iaas]
        resource: environment
  - name: another-job
    plan:
      - task: another-task
        tags: [some-iaas]
        resource: environment
  - name: destroy-environment-final
    plan:
      - get: environment
        resource: environment
        passed: [claim-environment]
      - task: destroy
        tags: [some-iaas]
YAML
  end

  let(:aws_extra_config) do
    <<YAML
- task: some-aws-task
  tags: aws
- task: another-aws-task
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

  describe '#add_aws_configure_tasks' do
    before do
      allow(File).to receive(:read).with('ci/pipelines/release/template/aws-extra-config')
        .and_return(aws_extra_config)
    end

    it 'adds all tasks from extra template file to configure-ert job' do
      task_adder.add_aws_configure_tasks(pipeline_yaml, 'aws-extra-config')

      expect(pipeline_yaml['jobs'][3]['plan'][1]['task']).to eq('some-aws-task')
      expect(pipeline_yaml['jobs'][3]['plan'][2]['task']).to eq('another-aws-task')
    end
  end

  describe '#add_vcloud_delete_installation_tasks' do
    before do
      allow(File).to receive(:read).with('ci/pipelines/release/template/vcloud-delete-installation.yml')
        .and_return(vcloud_extra_config)
    end

    it 'adds all tasks from extra template file to all destroy jobs' do
      task_adder.add_vcloud_delete_installation_tasks(pipeline_yaml)

      expect(pipeline_yaml['jobs'][0]['plan'][2]['task']).to eq('some-vcloud-task')
      expect(pipeline_yaml['jobs'][0]['plan'][3]['task']).to eq('another-vcloud-task')

      expect(pipeline_yaml['jobs'][5]['plan'][1]['task']).to eq('some-vcloud-task')
      expect(pipeline_yaml['jobs'][5]['plan'][2]['task']).to eq('another-vcloud-task')
    end
  end

  describe '#add_verify_internetless_job' do
    before do
      allow(File).to receive(:read).with('ci/pipelines/release/template/internetless-verification.yml')
        .and_return(verify_internetless_config)
    end

    it 'adds the verify_internetless job after deploy-ops-manager' do
      task_adder.add_verify_internetless_job(pipeline_yaml)

      expect(pipeline_yaml['jobs'][2]['name']).to eq('verify-internetless')
      expect(pipeline_yaml['jobs'][2]['plan'][0]['task']).to eq('verify')

      pipeline_yaml['jobs'][3]['plan'].each do |task|
        expect(task['passed']).to eq(['verify-internetless']) if task['passed']
      end
    end
  end
end
