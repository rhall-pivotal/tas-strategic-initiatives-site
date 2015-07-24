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
        trigger: true
        passed: [claim-environment]
      - get: environment
        trigger: true
        passed: [claim-environment]
      - task: destroy
        tags: [some-iaas]
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
        trigger: true
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

  describe '#add_aws_configure_tasks' do
    before do
      allow(File).to receive(:read).with('ci/pipelines/release/template/aws-extra-config')
        .and_return(aws_extra_config)
    end

    it 'adds all tasks from extra template file to configure-ert job' do
      task_adder.add_aws_configure_tasks(pipeline_yaml, 'aws-extra-config')

      expect(pipeline_yaml['jobs'][1]['plan'][1]['task']).to eq('some-aws-task')
      expect(pipeline_yaml['jobs'][1]['plan'][2]['task']).to eq('another-aws-task')
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

        expect(pipeline_yaml['jobs'][3]['plan'][1]['task']).to eq('some-vcloud-task')
        expect(pipeline_yaml['jobs'][3]['plan'][2]['task']).to eq('another-vcloud-task')
      end
    end
  end
end
