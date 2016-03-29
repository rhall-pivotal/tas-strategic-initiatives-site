require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.7.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
  end

  # This is needed until ops manager fixes automatic sizing bug: https://www.pivotaltracker.com/story/show/115004337
  it 'configures diego cell instance type' do
    instance_type = case env_settings['iaas_type']
                    when 'aws'
                      'm3.2xlarge'
                    when 'openstack'
                      'm1.xlarge'
                    when 'vsphere', 'vcloud'
                      '2xlarge.cpu'
                    end
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_vm_type_for_job('diego_cell', instance_type)
  end

  it 'configure postgres instance counts to zero' do
    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_instances_for_job('ccdb', 0)
    resource_config.set_instances_for_job('consoledb', 0)
    resource_config.set_instances_for_job('uaadb', 0)
  end

  it 'configures the availability zone' do
    if can_configure_availability_zones?('cf')
      unless any_availability_zones_have_been_selected_for_balancing?('cf')
        current_ops_manager.assign_azs_and_network_for_product(
          product_name: 'cf',
          zones: env_settings['ops_manager']['availability_zones'],
          network: env_settings['ops_manager']['networks'].first['name']
        )
      end
    end
  end

  it 'configures the domains' do
    domains_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('domains')
    domains_form.open_form
    domains_form.property('.cloud_controller.system_domain').set(elastic_runtime_settings['system_domain'])
    domains_form.property('.cloud_controller.apps_domain').set(elastic_runtime_settings['apps_domain'])
    domains_form.save_form
  end

  it 'configures networking' do
    case env_settings['iaas_type']
    when 'aws'
      configure_aws_load_balancers(elastic_runtime_settings)
    when 'vsphere', 'vcloud'
      configure_vsphere_ha_proxy(elastic_runtime_settings)
    when 'openstack'
      configure_openstack_ha_proxy(elastic_runtime_settings)
    end
  end

  it 'configures smtp' do
    if elastic_runtime_settings['smtp']
      smtp_form = current_ops_manager.product(elastic_runtime_settings['name']).product_form('smtp_config')
      smtp_form.open_form
      smtp_form.property('.properties.smtp_from').set(elastic_runtime_settings['smtp']['from'])
      smtp_form.property('.properties.smtp_address').set(elastic_runtime_settings['smtp']['address'])
      smtp_form.property('.properties.smtp_port').set(elastic_runtime_settings['smtp']['port'])
      smtp_form.nested_property('.properties.smtp_credentials', 'identity')
        .set(elastic_runtime_settings['smtp']['credentials']['identity'])
      smtp_form.nested_property('.properties.smtp_credentials', 'password')
        .set(elastic_runtime_settings['smtp']['credentials']['password'])
      smtp_form.property('.properties.smtp_enable_starttls_auto').set(elastic_runtime_settings['smtp']['enable_starttls_auto'])
      smtp_form.property('.properties.smtp_auth_mechanism').set(elastic_runtime_settings['smtp']['smtp_auth_mechanism'])
      smtp_form.save_form
    end
  end

  private

  def configure_vsphere_ha_proxy(elastic_runtime_settings)
    networking_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('networking')
    networking_form.open_form
    networking_form.property('.ha_proxy.static_ips').set(elastic_runtime_settings['ha_proxy_static_ips'])
    networking_form.fill_in_selector_property(
      selector_input_reference: '.properties.networking_point_of_entry',
      selector_name: '',
      selector_value: 'haproxy',
      sub_field_answers: {}
    )
    configure_ssl_cert(networking_form, elastic_runtime_settings, 'haproxy')
    networking_form.save_form
  end

  def configure_aws_load_balancers(elastic_runtime_settings)
    networking_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('networking')
    networking_form.open_form
    networking_form.fill_in_selector_property(
      selector_input_reference: '.properties.networking_point_of_entry',
      selector_name: '',
      selector_value: 'external_ssl',
      sub_field_answers: {}
    )
    configure_ssl_cert(networking_form, elastic_runtime_settings, 'external_ssl')
    networking_form.save_form

    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_instances_for_job('ha_proxy', 0)
    resource_config.set_elb_names_for_job('router', elastic_runtime_settings['elb_name'])
    resource_config.set_elb_names_for_job('diego_brain', elastic_runtime_settings['ssh_elb_name'])

    system_logging_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('syslog_aggregator')
    system_logging_form.open_form
    system_logging_form.property('.properties.logger_endpoint_port').set('4443')
    system_logging_form.save_form
  end

  def configure_openstack_ha_proxy(elastic_runtime_settings)
    networking_form =
      current_ops_manager.product(elastic_runtime_settings['name']).product_form('networking')
    networking_form.open_form
    networking_form.fill_in_selector_property(
      selector_input_reference: '.properties.networking_point_of_entry',
      selector_name: '',
      selector_value: 'haproxy',
      sub_field_answers: {}
    )
    configure_ssl_cert(networking_form, elastic_runtime_settings, 'haproxy')
    networking_form.save_form

    resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
    resource_config.set_floating_ips_for_job('ha_proxy', elastic_runtime_settings['ha_proxy_floating_ips'])
  end

  def configure_ssl_cert(networking_form, elastic_runtime_settings, selector_option)
    networking_form.property('.ha_proxy.skip_cert_verify').set(elastic_runtime_settings['trust_self_signed_certificates'])

    ssl_rsa_cert_property = '.properties.networking_point_of_entry][' + \
                            selector_option + '][.properties.networking_point_of_entry.' + \
                            selector_option + '.ssl_rsa_certificate'

    if elastic_runtime_settings['ssl_certificate']
      networking_form.nested_property(ssl_rsa_cert_property, 'cert_pem').set(elastic_runtime_settings['ssl_certificate'])
      networking_form.nested_property(ssl_rsa_cert_property, 'private_key_pem').set(elastic_runtime_settings['ssl_private_key'])
    else
      domain = elastic_runtime_settings['system_domain']
      networking_form.generate_self_signed_cert(
        "*.#{domain},*.login.#{domain},*.uaa.#{domain}",
        '.properties.networking_point_of_entry.' + selector_option + '.ssl_rsa_certificate',
        '.properties.networking_point_of_entry',
        selector_option
      )
    end
  end

  AVAILABILITY_ZONE_INPUT_SELECTOR = "input[name='product[availability_zone_references][]']"

  def availability_zones_for_product(product:)
    visit '/'
    click_on "show-#{product}-configure-action"

    az_selector = "show-#{product}-azs-and-network-assignment-action"

    click_on az_selector
    all("#{AVAILABILITY_ZONE_INPUT_SELECTOR}[checked='checked']").map do |checkbox|
      find("label[for='#{checkbox[:id]}']").text
    end
  end

  def any_availability_zones_have_been_selected_for_balancing?(product)
    availability_zones_for_product(product: product).length > 0
  end

  def can_configure_availability_zones?(product)
    visit '/'
    click_on "show-#{product}-configure-action"

    az_selector = "show-#{product}-azs-and-network-assignment-action"
    has_link? az_selector
  end
end
