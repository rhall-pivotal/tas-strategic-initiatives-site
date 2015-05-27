require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.5.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings.ops_manager.elastic_runtime }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings.ops_manager.username,
      password: env_settings.ops_manager.password,
    )
  end

  it 'configures the availability zone' do
    current_ops_manager.assign_availability_zones_for_product(product: 'cf', zones: env_settings.ops_manager.availability_zones)
  end

  it 'configures ha proxy or an ELB' do
    case env_settings.iaas_type
    when 'aws'
      configure_aws_load_balancers(elastic_runtime_settings)
    when 'vsphere'
      configure_vsphere_ha_proxy(elastic_runtime_settings)
    end
  end

  it 'configures ssl certificates' do
    security_config_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('security_config')
    security_config_form.open_form
    security_config_form.property('.ha_proxy.skip_cert_verify').set(elastic_runtime_settings.trust_self_signed_certificates)
    security_config_form.generate_self_signed_cert("*.#{elastic_runtime_settings.system_domain}")
    security_config_form.save_form
  end

  it 'configures cloud controller' do
    cloud_controller_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('cloud_controller')
    cloud_controller_form.open_form
    cloud_controller_form.property('.cloud_controller.system_domain').set(elastic_runtime_settings.system_domain)
    cloud_controller_form.property('.cloud_controller.apps_domain').set(elastic_runtime_settings.apps_domain)
    cloud_controller_form.save_form
  end

  it 'configures smtp' do
    smtp_form = current_ops_manager.product(elastic_runtime_settings.name).product_form('smtp_config')
    smtp_form.open_form
    smtp_form.property('.properties.smtp_from').set(elastic_runtime_settings.smtp.from)
    smtp_form.property('.properties.smtp_address').set(elastic_runtime_settings.smtp.address)
    smtp_form.property('.properties.smtp_port').set(elastic_runtime_settings.smtp.port)
    smtp_form.nested_property('.properties.smtp_credentials', 'identity').set(elastic_runtime_settings.smtp.credentials.identity)
    smtp_form.nested_property('.properties.smtp_credentials', 'password').set(elastic_runtime_settings.smtp.credentials.password)
    smtp_form.property('.properties.smtp_enable_starttls_auto').set(elastic_runtime_settings.smtp.enable_starttls_auto)
    smtp_form.property('.properties.smtp_auth_mechanism').set(elastic_runtime_settings.smtp.smtp_auth_mechanism)
    smtp_form.save_form
  end

  private

  def configure_vsphere_ha_proxy(elastic_runtime_settings)
    ips_and_ports_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('ha_proxy')
    ips_and_ports_form.open_form
    ips_and_ports_form.property('.ha_proxy.static_ips').set(elastic_runtime_settings.ha_proxy_static_ips)
    ips_and_ports_form.save_form
  end

  def configure_aws_load_balancers(elastic_runtime_settings)
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings.name)
    resource_config.set_instances_for_job('ha_proxy', 0)
    resource_config.set_elb_name_for_job('router', elastic_runtime_settings.elb_name)

    ips_and_ports_form =
      current_ops_manager.product(elastic_runtime_settings.name).product_form('ha_proxy')
    ips_and_ports_form.open_form
    ips_and_ports_form.property('.properties.logger_endpoint_port').set('4443')
    ips_and_ports_form.save_form
  end
end
