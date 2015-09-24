require 'spec_helper'
require 'tools/concourse_failure_reporter'
describe ConcourseFailureReporter do
  subject(:reporter) { ConcourseFailureReporter.new }

  let(:history) do
    [
      {
        'id' => 212,
        'name' => '1',
        'status' => 'succeeded',
        'job_name' => 'deploy-ops-manager-vsphere-clean',
        'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
      },
      {
        'id' => 211,
        'name' => '1',
        'status' => 'succeeded',
        'job_name' => 'destroy-environment-vsphere-upgrade',
        'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
      }
    ]
  end
  let(:new_build_info) do
    [
      {
        'id' => 212,

        'name' => '1',
        'status' => 'succeeded',
        'job_name' => 'deploy-ops-manager-vsphere-clean',
        'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
      },
      {
        'id' => 211,
        'name' => '1',
        'status' => 'succeeded',
        'job_name' => 'destroy-environment-vsphere-upgrade',
        'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
      }
    ]
  end

  before do
    allow(Slack::Post).to(receive(:configure))
    allow(File).to(receive(:write))
    allow(File).to receive(:read).with('/Users/pivotal/reporter_history.json').and_return(history.to_json)
    allow(YAML).to(
      receive(:load_file)
        .with('/Users/pivotal/.flyrc')
        .and_return(
          'targets' => {
            'ci' => {
              'api' => 'api',
              'username' => 'username',
              'password' => 'password',
            }
          }
        )
    )
  end

  it 'registers with slack' do
    expect(Slack::Post).to(receive(:configure))
    reporter
  end

  context 'when there are no new builds' do
    let(:new_build_info) { history }
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end
    it 'does not post a message' do
      reporter.run
    end
  end
  context 'when there is a new build that succeeded' do
    let(:new_build_info) do
      [
        {
          'id' => 213,
          'name' => '2',
          'status' => 'succeeded',
          'job_name' => 'deploy-ops-manager-vsphere-clean',
          'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/2'
        },
        {
          'id' => 212,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'deploy-ops-manager-vsphere-clean',
          'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
        },
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end
    it 'does not post a message' do
      reporter.run
    end
  end
  context 'when an old build went from pending to succeeded' do
    let(:new_build_info) do
      [
        {
          'id' => 212,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'deploy-ops-manager-vsphere-clean',
          'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
        },
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    let(:history) do
      [
        {
          'id' => 212,
          'name' => '1',
          'status' => 'pending',
          'job_name' => 'deploy-ops-manager-vsphere-clean',
          'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
        },
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end
    it 'does not post a message' do
      reporter.run
    end
  end
  context 'when an new build is present and pending' do
    let(:new_build_info) do
      [
        {
          'id' => 212,
          'name' => '1',
          'status' => 'pending',
          'job_name' => 'deploy-ops-manager-vsphere-clean',
          'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1'
        },
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    let(:history) do
      [
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end
    it 'does not post a message' do
      reporter.run
    end
  end
  context 'when an new build is present and failed' do
    let(:new_build_info) do
      [
        failed_job,
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    let(:failed_job) do
      { 'id' => 212,
        'name' => '1',
        'status' => 'failed',
        'job_name' => 'deploy-ops-manager-vsphere-clean',
        'url' => '/pipelines/ert-1.5/jobs/deploy-ops-manager-vsphere-clean/builds/1' }
    end

    let(:failed_job_message) do
      <<"MESSAGE"
Failed build #{failed_job['id']}
Job: #{failed_job['job_name']}
URL: api#{failed_job['url']}
MESSAGE
    end

    let(:history) do
      [
        {
          'id' => 211,
          'name' => '1',
          'status' => 'succeeded',
          'job_name' => 'destroy-environment-vsphere-upgrade',
          'url' => '/pipelines/ert-1.5/jobs/destroy-environment-vsphere-upgrade/builds/1'
        }
      ]
    end
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end
    it 'posts a message' do
      expect(Slack::Post).to(receive(:post).with(failed_job_message))
      reporter.run
    end
  end

  context 'saving history' do
    before do
      allow(reporter).to receive(:new_build_info).and_return(new_build_info)
    end

    it 'writes the new build info to the history file' do
      expect(File).to receive(:write) do |filename, contents|
        expect(filename).to eq('/Users/pivotal/reporter_history.json')
        expect(contents).to eq(new_build_info.to_json)
      end

      reporter.run
    end

    it 'reads the history file for the old values' do
      expect(File).to receive(:read).with('/Users/pivotal/reporter_history.json').and_return(history.to_json)
      expect(reporter.history).to eq(history)
    end

    context 'when the history file does not exist' do
      before do
        expect(File).to receive(:read).with('/Users/pivotal/reporter_history.json').and_raise(Errno::ENOENT.new)
      end
      it 'does not report any failures' do
        reporter.run
      end
      it 'writes the new build info to the history file' do
        expect(File).to receive(:write) do |filename, contents|
          expect(filename).to eq('/Users/pivotal/reporter_history.json')
          expect(contents).to eq(new_build_info.to_json)
        end
        reporter.run
      end
    end
  end
end
