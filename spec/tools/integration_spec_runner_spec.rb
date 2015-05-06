require 'tools/integration_spec_runner'

RSpec.describe 'IntegrationSpecRunner' do
  subject(:integration_spec_runner) do
    IntegrationSpecRunner.new(
      environment: environment,
      ert_version: ert_version,
      om_version: om_version
    )
  end

  let(:environment) { nil }
  let(:ert_version) { nil }
  let(:om_version) { nil }

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
      let(:ert_version) { 'UNSUPPORTED' }

      it 'raises a helpful error' do
        expect do
          integration_spec_runner.configure_ert
        end.to raise_error(IntegrationSpecRunner::UnsupportedErtVersion, 'Version "UNSUPPORTED" is not supported')
      end
    end

    context 'when ert_version is not passed' do
      let(:ert_version) { nil }

      it 'does not raise an error' do
        expect do
          IntegrationSpecRunner.new(
            environment: 'foo',
            om_version: '1.5',
          )
        end.not_to raise_error
      end
    end
  end

  %w(1.4 1.5).each do |version|
    describe "#configure_ert #{version}" do
      let(:ert_version) { version }

      it 'runs the correct version of configure ert' do
        expect(RSpec::Core::Runner).to receive(:run).with(["integration/ERT-#{ert_version}/configure_ert_spec.rb"])

        integration_spec_runner.configure_ert
      end
    end
  end
end
