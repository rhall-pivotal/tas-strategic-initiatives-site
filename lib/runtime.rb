# rubocop:disable Metrics/LineLength

require 'opsmgr/log'
require 'opsmgr/configurator'
require 'opsmgr/settings/microbosh/network'
require 'tools/self_signed_rsa_certificate'
require 'opsmgr/environment'

class Runtime
  include Opsmgr::Loggable

  PRODUCT_NAME = 'cf'.freeze

  def self.build(environment)
    new(Opsmgr::Configurator.build(environment))
  end

  def initialize(ops_manager_configurator)
    @ops_manager_configurator = ops_manager_configurator
    @environment = ops_manager_configurator.environment
  end

  # rubocop:disable Metrics/AbcSize, Metrics/MethodLength, Metrics/CyclomaticComplexity
  def configure
    log.info "configuring #{environment.name}'s Elastic Runtime"

    ops_manager_configurator.add_product(PRODUCT_NAME)
    ops_manager_configurator.configure do |settings|
      runtime = settings.product(PRODUCT_NAME)
      runtime.disabled_post_deploy_errand_names(environment.ers_configuration[:disabled_post_deploy_errand_names])
      runtime.singleton_availability_zone = settings.default_availability_zone_guid
      runtime.availability_zone_references = [settings.default_availability_zone_guid]
      runtime.network_reference = settings.network_guid('default')

      runtime.jobs.each do |job|
        job.networks = [settings.network_guid('default')]
      end

      environment.ers_configuration[:jobs].each do |job_name, job_config|
        runtime.for_job(job_name.to_s) do |job|
          job.resource('instances').value = job_config[:instances] if job_config[:instances]
        end
      end

      runtime.for_job('ha_proxy') do |job|
        job.property('static_ips').value = environment.ers_configuration.fetch(:ha_proxy_ips, []).join(',')
        job.property('ssl_rsa_certificate').value = ha_proxy_ssl_rsa_certificate_value

        skip_cert_verify = job.property('skip_cert_verify')
        skip_cert_verify.value = environment.ers_configuration.fetch(:trust_self_signed_certificates) unless skip_cert_verify.nil?
      end

      runtime.for_job('router') do |job|
        job.elb_name = environment.ers_configuration[:elb_name]
      end

      runtime.for_job('cloud_controller') do |job|
        job.property('system_domain').value = environment.ers_configuration[:system_domain]
        job.property('apps_domain').value = environment.ers_configuration[:apps_domain]
        job.property('max_file_size').value = 1024 if job.property('max_file_size')
      end

      if environment.ers_configuration.key?(:smtp) && !runtime.property('smtp_from').nil?
        runtime.set_property('smtp_from', environment.ers_configuration[:smtp][:from])
        runtime.set_property('smtp_address', environment.ers_configuration[:smtp][:address])
        runtime.set_property('smtp_port', environment.ers_configuration[:smtp][:port])
        runtime.set_property('smtp_credentials',
                             'identity' => environment.ers_configuration[:smtp][:credentials][:identity],
                             'password' => environment.ers_configuration[:smtp][:credentials][:password])

        runtime.set_property('smtp_enable_starttls_auto', environment.ers_configuration[:smtp][:enable_starttls_auto])
        runtime.set_property('smtp_auth_mechanism', environment.ers_configuration[:smtp][:auth_mechanism])
      end

      if environment.ers_configuration[:logging_port]
        runtime.set_property('logger_endpoint_port', environment.ers_configuration[:logging_port])
      end
    end
  end
  # rubocop:enable Metrics/AbcSize, Metrics/MethodLength, Metrics/CyclomaticComplexity

  private

  attr_reader :ops_manager_configurator, :environment

  # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
  def ha_proxy_ssl_rsa_certificate_value
    if environment.ers_configuration[:ssl_certificate]
      {
        'cert_pem' => environment.ers_configuration[:ssl_certificate],
        'private_key_pem' => environment.ers_configuration[:ssl_private_key]
      }
    else
      certificate = Tools::SelfSignedRsaCertificate.generate(environment.ers_configuration[:ssl_cert_domains].split(','))
      {
        'cert_pem' => certificate.cert_pem,
        'private_key_pem' => certificate.private_key_pem
      }
    end
  end
  # rubocop:enable Metrics/AbcSize, Metrics/MethodLength
end
