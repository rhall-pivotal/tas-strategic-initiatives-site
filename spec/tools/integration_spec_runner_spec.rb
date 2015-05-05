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

  describe '#upload_ert' do
    let(:ert_version) { '1.5' }
    let(:product_path) { '/some/path/to.pivotal' }

    it 'runs the correct version of configure ert' do
      expect(RSpec::Core::Runner).to receive(:run).with(["integration/ERT-#{ert_version}/upload_ert_spec.rb"])
      expect(ENV).to receive(:[]=).with('PRODUCT_PATH', product_path)

      integration_spec_runner.upload_ert(product_path)
    end
  end

  describe '#configure_ert' do
    let(:ert_version) { '1.5' }

    it 'runs the correct version of configure ert' do
      expect(RSpec::Core::Runner).to receive(:run).with(["integration/ERT-#{ert_version}/configure_ert_spec.rb"])

      integration_spec_runner.configure_ert
    end
  end

  describe '#configure_microbosh' do
    let(:om_version) { '1.5' }

    it 'configures microbosh' do
      expect(RSpec::Core::Runner).to receive(:run).with(['integration/microbosh/configure_microbosh_spec.rb'])

      integration_spec_runner.configure_microbosh
    end
  end

  describe '#install_microbosh' do
    let(:om_version) { '1.5' }

    it 'install microbosh' do
      expect(RSpec::Core::Runner).to receive(:run).with(['integration/microbosh/install_microbosh_spec.rb'])

      integration_spec_runner.install_microbosh
    end
  end
end
