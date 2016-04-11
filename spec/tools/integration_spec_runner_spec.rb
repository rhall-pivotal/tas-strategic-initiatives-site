require 'tools/integration_spec_runner'

RSpec.describe 'IntegrationSpecRunner' do
  subject(:integration_spec_runner) do
    IntegrationSpecRunner.new(
      environment: environment,
      ert_version: ert_version,
      om_version: om_version
    )
  end

  let(:environment) { 'some-env' }
  let(:om_version) { 'some-version' }

  before { allow(ENV).to receive(:[]=) }

  describe '#intialize' do
    context 'with valid parameters' do
      it 'does not raise an error' do
        expect do
          IntegrationSpecRunner.new(
            environment: 'foo',
            om_version: '1.5',
            ert_version: '1.5',
          )
        end.not_to raise_error
      end

      it 'sets the ENV[ENVIRONMENT_NAME]' do
        expect(ENV).to receive(:[]=).with('ENVIRONMENT_NAME', 'foo')
        IntegrationSpecRunner.new(
          environment: 'foo',
          om_version: '1.5',
          ert_version: '1.5',
        )
      end

      it 'sets the ENV[OM_VERSION]' do
        expect(ENV).to receive(:[]=).with('OM_VERSION', '1.5')
        IntegrationSpecRunner.new(
          environment: 'foo',
          om_version: '1.5',
          ert_version: '1.5',
        )
      end
    end

    context 'with an unsupported ert_version' do
      it 'raises a helpful error' do
        expect do
          IntegrationSpecRunner.new(
            environment: 'foo',
            om_version: '1.5',
            ert_version: 'UNSUPPORTED',
          )
        end.to raise_error(IntegrationSpecRunner::UnsupportedErtVersion, 'Version "UNSUPPORTED" is not supported')
      end
    end
  end

  context 'when using a 1.6 specific task' do
    let(:ert_version) { '1.6' }

    describe '#configure_postgres' do
      it 'runs the correct version of configure postgres' do
        expect(RSpecExiter).to receive(:exit_rspec).with(0)
        expect(RSpec::Core::Runner).to receive(:run).with(
          ['integration/ERT-1.6/configure_postgres_spec.rb']
        ).and_return(0)

        integration_spec_runner.configure_postgres
      end

      context 'when called with a unsupported version' do
        let(:ert_version) { '1.7' }
        it 'returns an UnsupportedErtVersion error' do
          expect do
            integration_spec_runner.configure_postgres
          end.to raise_error(IntegrationSpecRunner::UnsupportedErtVersion, 'Version 1.7 is not supported for this task')
        end
      end
    end
  end

  %w(1.5 1.6).each do |version|
    describe 'configuring ert' do
      let(:ert_version) { version }

      describe "#configure_ert #{version}" do
        it 'runs the correct version of configure ert' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to receive(:run).with(
            ["integration/ERT-#{ert_version}/configure_ert_spec.rb"]
          ).and_return(0)

          integration_spec_runner.configure_ert
        end
      end

      describe "#configure_external_dbs #{version}" do
        it 'runs the correct version of configure external dbs' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to receive(:run).with(
            ["integration/ERT-#{ert_version}/configure_external_dbs_spec.rb"]
          ).and_return(0)

          integration_spec_runner.configure_external_dbs
        end
      end

      describe "#configure_external_file_storage #{version}" do
        it 'runs the correct version of configure external file storage' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to receive(:run).with(
            ["integration/ERT-#{ert_version}/configure_external_file_storage_spec.rb"]
          ).and_return(0)

          integration_spec_runner.configure_external_file_storage
        end
      end

      describe "#configure_multi_az_instance_counts #{version}" do
        it 'runs the correct version of configure_multi_az_instance_counts' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to receive(:run).with(
            ["integration/ERT-#{ert_version}/configure_multi_az_instance_counts_spec.rb"]
          ).and_return(0)

          integration_spec_runner.configure_multi_az_instance_counts
        end
      end

      describe "#configure_experimental_features #{version}" do
        it 'runs the correct version of configure experimental features' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to(
            receive(:run).with(
              ["integration/ERT-#{ert_version}/configure_experimental_features_spec.rb"]
            ).and_return(0)
          )

          integration_spec_runner.configure_experimental_features
        end
      end

      describe "#disable_http_traffic #{version}" do
        it 'runs the correct version of configure security features' do
          expect(RSpecExiter).to receive(:exit_rspec).with(0)
          expect(RSpec::Core::Runner).to(
            receive(:run).with(
              ["integration/ERT-#{ert_version}/disable_http_traffic.rb"]
            ).and_return(0)
          )

          integration_spec_runner.disable_http_traffic
        end
      end
    end
  end
end
