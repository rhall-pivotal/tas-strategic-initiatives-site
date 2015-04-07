require 'spec_helper'

require 'runtime'

describe Runtime, :teapot do
  include DiffHelpers

  let(:settings) do
    {
      environments: {
        test: {
          web_auth: {
            web_user: 'test_web_user',
            web_password: 'test_web_password'
          },
          vcenter: {
            host: '127.1.1.1',
            user: 'vcenter_user',
            password: 'vcenter_password'
          },
          location: {
            datacenter: 'test-dc',
            cluster: 'test-cl',
            datastore: 'test-ds',
            network: 'Test Network',
            folder: 'test'
          },
          ops_manager_settings: {
            ip: '127.0.0.1',
            port: '8741',
            protocol: 'http',
            netmask: '255.255.254.0',
            gateway: '127.0.0.254',
            reserved_ip_ranges: '127.0.0.1-127.0.0.16,127.0.0.100-127.0.1.255',
            dns: '192.168.2.3',
            ntp_servers: ['time.example.com'],
            vm_password: 'test_vm_password' },
          ers_configuration: {
            trust_self_signed_certificates: false,
            elb_name: 'my-elb',
            ssl_cert_domains: '*.ssl.example.com',
            system_domain: 'system.example.com',
            apps_domain: 'apps.example.com',
            disabled_post_deploy_errand_names: %w(push-console smoke-tests),
            smtp: {
              from: 'reply_to@example.com',
              address: 'smtp.example.com',
              port: 587,
              credentials: {
                identity: 'notifications_id',
                password: 'secret'
              },
              enable_starttls_auto: true,
              auth_mechanism: 'none' },
            jobs: {
              clock_global: {
                instances: 0
              }
            }
          }
        }
      }
    }
  end

  let(:test_env) { Opsmgr::Environment.build(:test, settings) }
  let(:test_user) { test_env.web_auth.user }
  let(:api_client) { Opsmgr::Api::Client.new(test_env) }
  let(:ops_manager_configurator) { Opsmgr::Configurator.new(test_env, api_client) }
  let(:cf_details) do
    req = Net::HTTP::Get.new('/teapot/component_details/cf')
    JSON.parse(teapot_client.request(req).body)
  end

  subject(:runtime) { Runtime.new(ops_manager_configurator) }

  def property_value(job_name, property_name)
    property(job_name, property_name)['value']
  end

  def resource_value(job_name, resource_name)
    resource(job_name, resource_name).value
  end

  def top_level_property_value(property_name)
    top_level_property(property_name)['value']
  end

  def job(job_name)
    Opsmgr::Settings::Microbosh::JobList.new(cf_details['jobs']).find { |j| j.name == job_name }
  end

  def property(job_name, property_name)
    Opsmgr::Settings::Microbosh::PropertyList.new(job(job_name)['properties']).find { |p| p.name == property_name }
  end

  def top_level_property(property_name)
    Opsmgr::Settings::Microbosh::Product.new(cf_details).property(property_name)
  end

  def resource(job_name, resource_name)
    Opsmgr::Settings::Microbosh::Job.new(job(job_name)).resource(resource_name)
  end

  describe '.build' do
    it 'correctly builds a Runtime product' do
      set_teapot_version('1.3')

      expect(Runtime.build(test_env)).to be_a(Runtime)
    end
  end

  context 'when cf version is 1.2' do
    before { set_teapot_version('1.2') }

    describe 'configuring HAProxy and Cloud Controller' do
      let(:expected_skip_cert_verify) { false }
      before do
        Tools::SelfSignedRsaCertificate.should_receive(:generate)
          .with(['*.ssl.example.com'])
          .and_return(double('certificate', private_key_pem: 'generated_private_key_pem', cert_pem: 'generated_cert_pem'))
      end

      shared_examples 'configures and installs the CF tile' do
        it 'enters our configuration and appears in the installation.yml' do
          runtime.configure

          expect(property_value('ha_proxy', 'skip_cert_verify')).to eq(expected_skip_cert_verify)
          diff_assert(
            property_value('ha_proxy', 'ssl_rsa_certificate'),
            'private_key_pem' => 'generated_private_key_pem', 'cert_pem' => 'generated_cert_pem'
          )
          expect(property_value('cloud_controller', 'system_domain')).to eq('system.example.com')
          expect(property_value('cloud_controller', 'apps_domain')).to eq('apps.example.com')
          expect(property_value('cloud_controller', 'max_file_size')).to eq(1024)
        end

        it 'sets a default network on each job' do
          runtime.configure

          expect(job('ha_proxy')['network_references']).to eq(['guid-for-the-default-network'])
          expect(job('cloud_controller')['network_references']).to eq(['guid-for-the-default-network'])
        end

        it 'does not configure notifications' do
          runtime.configure

          expect(top_level_property('smtp_from')).to be_nil
        end
      end

      context 'when the CF tile has not yet been added' do
        include_examples 'configures and installs the CF tile'
      end

      context 'when the CF tile has already been added' do
        before do
          add_teapot_product('cf', '1.2.0.0')
        end

        include_examples 'configures and installs the CF tile'
      end

      context 'when the installation requires more cores than the environment has' do
        before do
          enable_number_of_cores_error
        end

        include_examples 'configures and installs the CF tile'
      end

      context 'when the environment is configured to trust self-signed certificates' do
        let(:expected_skip_cert_verify) { true }

        before do
          actual_ers_config = test_env.ers_configuration
          fake_ers_config = actual_ers_config.merge(trust_self_signed_certificates: true)
          expect(test_env).to receive(:ers_configuration).and_return(fake_ers_config).at_least(1).times
        end

        include_examples 'configures and installs the CF tile'
      end

      context 'when the user wants to skip errands' do
        it 'should set the "disabled_post_deploy_errand_names"' do
          actual_ers_config = test_env.ers_configuration
          fake_ers_config = actual_ers_config.merge(disabled_post_deploy_errand_names: %w(a c))
          expect(test_env).to receive(:ers_configuration).and_return(fake_ers_config).at_least(1).times
          runtime.configure

          cf_details = begin
            req = Net::HTTP::Get.new('/teapot/component_details/cf')
            JSON.parse(teapot_client.request(req).body)
          end

          expect(cf_details['disabled_post_deploy_errand_names']).to eq(%w(a c))
        end
      end
    end
  end

  context 'when cf version is 1.3' do
    before { set_teapot_version('1.3') }

    describe 'configuring CF jobs' do
      let(:expected_skip_cert_verify) { false }

      before do
        Tools::SelfSignedRsaCertificate.should_receive(:generate)
          .with(['*.ssl.example.com'])
          .and_return(double('certificate', private_key_pem: 'generated_private_key_pem', cert_pem: 'generated_cert_pem'))
      end

      context 'when the CF tile has not yet been added' do
        it 'enters our configuration and appears in the installation.yml' do
          runtime.configure

          expect(property_value('ha_proxy', 'skip_cert_verify')).to eq(expected_skip_cert_verify)
          diff_assert(
            property_value('ha_proxy', 'ssl_rsa_certificate'),
            'private_key_pem' => 'generated_private_key_pem', 'cert_pem' => 'generated_cert_pem'
          )
          expect(property_value('cloud_controller', 'system_domain')).to eq('system.example.com')
          expect(property_value('cloud_controller', 'apps_domain')).to eq('apps.example.com')
          expect(property_value('cloud_controller', 'max_file_size')).to eq(1024)
        end

        it 'fills in notifications properties when smtp configuration is provided' do
          allow(ENV).to receive(:[]).and_call_original
          allow(ENV).to receive(:[]).with('REL_ENG_TEST_SMTP_PASSWORD') { 'secret' }

          runtime.configure

          expect(top_level_property_value('smtp_from')).to eq('reply_to@example.com')
          expect(top_level_property_value('smtp_address')).to eq('smtp.example.com')
          expect(top_level_property_value('smtp_port')).to eq(587)
          expect(top_level_property_value('smtp_credentials')).to eq('identity' => 'notifications_id', 'password' => 'secret')
          expect(top_level_property_value('smtp_enable_starttls_auto')).to be true
          expect(top_level_property_value('smtp_auth_mechanism')).to eq('none')
        end

        it 'sets a default network on each job' do
          runtime.configure

          expect(job('ha_proxy')['network_references']).to eq(['guid-for-the-default-network'])
          expect(job('cloud_controller')['network_references']).to eq(['guid-for-the-default-network'])
        end

        it 'assigns a singleton availability zone' do
          runtime.configure

          expect(cf_details['singleton_availability_zone_reference']).to eq('guid-for-the-availability-zone')
        end

        it 'assigns the availability zone references' do
          runtime.configure

          expect(cf_details['availability_zone_references']).to eq(['guid-for-the-availability-zone'])
        end

        it 'assigns the network_reference' do
          runtime.configure

          expect(cf_details['network_reference']).to eq('guid-for-the-default-network')
        end
      end
    end
  end

  context 'when cf version is 1.4' do
    before { set_teapot_version('1.4') }

    describe 'configuring CF jobs' do
      let(:expected_skip_cert_verify) { false }

      before do
        allow(Tools::SelfSignedRsaCertificate).to receive(:generate)
          .with(['*.ssl.example.com'])
          .and_return(double('certificate', private_key_pem: 'generated_private_key_pem', cert_pem: 'generated_cert_pem'))
      end

      context 'when the CF tile has not yet been added' do
        context 'when ha_proxy_ips are set' do
          before do
            settings[:environments][:test][:ers_configuration][:ha_proxy_ips] = ['192.168.2.4']
          end
          it 'sets them correctly' do
            runtime.configure

            expect(property_value('ha_proxy', 'static_ips')).to eq('192.168.2.4')
          end
        end

        context 'when ENV[INTERNET] is false' do
          before do
            allow(ENV).to receive(:[]).with(anything)
            allow(ENV).to receive(:[]).with('INTERNET').and_return('false')
          end

          it 'sets acceptance_tests.internet_available to false' do
            runtime.configure

            expect(property_value('acceptance-tests', 'internet_available')).to eq(false)
          end
        end

        it 'enters our configuration and appears in the installation.yml' do
          runtime.configure

          expect(property_value('ha_proxy', 'skip_cert_verify')).to eq(expected_skip_cert_verify)
          diff_assert(
            property_value('ha_proxy', 'ssl_rsa_certificate'),
            'private_key_pem' => 'generated_private_key_pem', 'cert_pem' => 'generated_cert_pem'
          )
          expect(property_value('cloud_controller', 'system_domain')).to eq('system.example.com')
          expect(property_value('cloud_controller', 'apps_domain')).to eq('apps.example.com')
          expect(property_value('cloud_controller', 'max_file_size')).to eq(1024)
        end

        it 'uses the provided certificate and private key' do
          settings[:environments][:test][:ers_configuration][:ssl_certificate] = 'provided_ssl_cert'
          settings[:environments][:test][:ers_configuration][:ssl_private_key] = 'provided_ssl_private_key'

          runtime.configure

          diff_assert(
            property_value('ha_proxy', 'ssl_rsa_certificate'),
            'private_key_pem' => 'provided_ssl_private_key', 'cert_pem' => 'provided_ssl_cert'
          )
        end

        it 'fills in notifications properties when smtp configuration is provided' do
          allow(ENV).to receive(:[]).and_call_original
          allow(ENV).to receive(:[]).with('REL_ENG_TEST_SMTP_PASSWORD') { 'secret' }

          runtime.configure

          expect(top_level_property_value('smtp_from')).to eq('reply_to@example.com')
          expect(top_level_property_value('smtp_address')).to eq('smtp.example.com')
          expect(top_level_property_value('smtp_port')).to eq(587)
          expect(top_level_property_value('smtp_credentials')).to eq('identity' => 'notifications_id', 'password' => 'secret')
          expect(top_level_property_value('smtp_enable_starttls_auto')).to be true
          expect(top_level_property_value('smtp_auth_mechanism')).to eq('none')
        end

        it 'sets a default network on each job' do
          runtime.configure

          expect(job('ha_proxy')['network_references']).to eq(['guid-for-the-default-network'])
          expect(job('cloud_controller')['network_references']).to eq(['guid-for-the-default-network'])
        end

        it 'assigns a singleton availability zone' do
          runtime.configure

          expect(cf_details['singleton_availability_zone_reference']).to eq('guid-for-the-availability-zone')
        end

        it 'assigns the availability zone references' do
          runtime.configure

          expect(cf_details['availability_zone_references']).to eq(['guid-for-the-availability-zone'])
        end

        it 'assigns the network_reference' do
          runtime.configure

          expect(cf_details['network_reference']).to eq('guid-for-the-default-network')
        end

        it 'sets the elb_name' do
          runtime.configure

          expect(job('router')['elb_name']).to eq('my-elb')
        end

        it 'sets job instance counts' do
          runtime.configure

          expect(job('clock_global')['instances'].first['value']).to eq(0)
        end

        it 'sets the logger_endpoint_port' do
          settings[:environments][:test][:ers_configuration][:logging_port] = 1234

          runtime.configure

          expect(top_level_property_value('logger_endpoint_port')).to eq(1234)
        end
      end
    end
  end
end
