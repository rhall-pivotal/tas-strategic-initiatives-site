require 'opsmgr/ui_helpers/config_helper'

RSpec.describe 'Configure Elastic Runtime 1.8.X', order: :defined do
  let(:current_ops_manager) { ops_manager_driver }
  let(:env_settings) { fetch_environment_settings }

  let(:elastic_runtime_settings) { env_settings['ops_manager']['elastic_runtime'] }

  it 'logs in' do
    current_ops_manager.setup_page.setup_or_login(
      user: env_settings['ops_manager']['username'],
      password: env_settings['ops_manager']['password'],
    )
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
    networking_form = current_ops_manager.product(elastic_runtime_settings['name']).product_form('networking')

    case env_settings['iaas_type']
    when 'aws'
    networking_form.open_form
      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.tcp_routing',
        selector_name: 'enable',
        selector_value: 'enable',
        sub_field_answers: {
          '.properties.tcp_routing.enable.reservable_ports' => {
            attribute_value: '1024-1123',
          },
        },
      )

      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.networking_point_of_entry',
        selector_name: '',
        selector_value: 'external_ssl',
        sub_field_answers: {}
      )
      networking_form.property('.ha_proxy.skip_cert_verify').set(elastic_runtime_settings['trust_self_signed_certificates'])
      ssl_rsa_cert_property = '.properties.networking_point_of_entry][external_ssl][.properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate'
      if elastic_runtime_settings['ssl_certificate']
        networking_form.nested_property(ssl_rsa_cert_property, 'cert_pem').set(elastic_runtime_settings['ssl_certificate'])
        networking_form.nested_property(ssl_rsa_cert_property, 'private_key_pem').set(elastic_runtime_settings['ssl_private_key'])
      else
        domain = elastic_runtime_settings['system_domain']
        networking_form.generate_self_signed_cert(
          "*.#{domain},*.login.#{domain},*.uaa.#{domain}",
          '.properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate',
          '.properties.networking_point_of_entry',
          'external_ssl'
        )
      end

      networking_form.property('.properties.logger_endpoint_port').set('4443')
      networking_form.save_form

      resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
      resource_config.set_elb_names_for_job('router', elastic_runtime_settings['elb_name'])
      resource_config.set_elb_names_for_job('diego_brain', elastic_runtime_settings['ssh_elb_name'])
      resource_config.set_elb_names_for_job('tcp_router', elastic_runtime_settings['tcp_elb_name'])

    when 'vsphere', 'vcloud'
      networking_form.open_form
      networking_form.property('.ha_proxy.static_ips').set(elastic_runtime_settings['ha_proxy_static_ips'])
      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.networking_point_of_entry',
        selector_name: '',
        selector_value: 'haproxy',
        sub_field_answers: {}
      )

      networking_form.property('.ha_proxy.skip_cert_verify').set(elastic_runtime_settings['trust_self_signed_certificates'])
      ssl_rsa_cert_property = '.properties.networking_point_of_entry][haproxy][.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate'
      if elastic_runtime_settings['ssl_certificate']
        networking_form.nested_property(ssl_rsa_cert_property, 'cert_pem').set(elastic_runtime_settings['ssl_certificate'])
        networking_form.nested_property(ssl_rsa_cert_property, 'private_key_pem').set(elastic_runtime_settings['ssl_private_key'])
      else
        domain = elastic_runtime_settings['system_domain']
        networking_form.generate_self_signed_cert(
          "*.#{domain},*.login.#{domain},*.uaa.#{domain}",
          '.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate',
          '.properties.networking_point_of_entry',
          'haproxy'
        )
      end

      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.tcp_routing',
        selector_name: 'enable',
        selector_value: 'enable',
        sub_field_answers: {
          '.properties.tcp_routing.enable.reservable_ports' => {
            attribute_value: '1024-1123',
          },
        },
      )
      networking_form.property('.tcp_router.static_ips').set(elastic_runtime_settings['tcp_router_static_ips'])
      networking_form.save_form
    when 'openstack'
      networking_form.open_form
      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.tcp_routing',
        selector_name: 'enable',
        selector_value: 'enable',
        sub_field_answers: {
          '.properties.tcp_routing.enable.reservable_ports' => {
            attribute_value: '1024-1123',
          },
        },
      )

      networking_form.fill_in_selector_property(
        selector_input_reference: '.properties.networking_point_of_entry',
        selector_name: '',
        selector_value: 'haproxy',
        sub_field_answers: {}
      )

      networking_form.property('.ha_proxy.skip_cert_verify').set(elastic_runtime_settings['trust_self_signed_certificates'])
      ssl_rsa_cert_property = '.properties.networking_point_of_entry][haproxy][.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate'
      if elastic_runtime_settings['ssl_certificate']
        networking_form.nested_property(ssl_rsa_cert_property, 'cert_pem').set(elastic_runtime_settings['ssl_certificate'])
        networking_form.nested_property(ssl_rsa_cert_property, 'private_key_pem').set(elastic_runtime_settings['ssl_private_key'])
      else
        domain = elastic_runtime_settings['system_domain']
        networking_form.generate_self_signed_cert(
          "*.#{domain},*.login.#{domain},*.uaa.#{domain}",
          '.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate',
          '.properties.networking_point_of_entry',
          'haproxy'
        )
      end
      networking_form.save_form

      resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
      resource_config.set_floating_ips_for_job('ha_proxy', elastic_runtime_settings['ha_proxy_floating_ips'])

      resource_config = current_ops_manager.product_resources_configuration(elastic_runtime_settings['name'])
      resource_config.set_floating_ips_for_job('tcp_router', elastic_runtime_settings['tcp_router_floating_ips'])
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

  it 'acknowledges the security text box' do
    if can_acknowledge_security?(elastic_runtime_settings['name'])
      security_form =
        current_ops_manager.product(elastic_runtime_settings['name']).product_form('application_security_groups')
      security_form.open_form
      security_form.property('.properties.security_acknowledgement').set('I agree')
      security_form.save_form
    end
  end

  private

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

  def can_acknowledge_security?(product)
    visit '/'
    click_on "show-#{product}-configure-action"
    has_link? "show-application_security_groups-action"
  end
end
