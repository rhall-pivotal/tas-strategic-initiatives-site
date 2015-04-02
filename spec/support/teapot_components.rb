# rubocop:disable Metrics/ClassLength
# rubocop:disable Metrics/MethodLength
class TeapotComponents
  def names_for_installation
    ['cf']
  end

  def names_for_installation_settings
    ['cf']
  end

  def components_for(version)
    components = []
    case version
    when '1.1'
      components << cf_1_1
    when '1.2'
      components << cf_1_2
    when '1.3', '1.4' # FIXME: handle identifier #88366870 <- this is not truly 1.4 compatible
      components << cf_1_3
    else
      fail 'unknown version'
    end
    components
  end

  private

  def cf_1_1
    {
      name: 'cf',
      status: 'available',
      config: {
        'product_version' => '1.1.0.0',
        'type' => 'cf',
        'jobs' => [{
          'type' => 'ha_proxy',
          'properties' => [{
            'definition' => 'static_ips'
          }, {
            'definition' => 'ssl_rsa_certificate'
          }]
        }, {
          'type' => 'cloud_controller',
          'properties' => [{
            'definition' => 'system_domain'
          }, {
            'definition' => 'apps_domain'
          }]
        }]
      },
      product_settings_endpoint: '/components/cf/forms/ha_proxy/edit'
    }
  end

  def cf_1_2
    {
      name: 'cf',
      status: 'available',
      config: {
        'product_version' => '1.2.0.0',
        'type' => 'cf',
        'jobs' => [{
          'type' => 'ha_proxy',
          'properties' => [{
            'definition' => 'static_ips'
          }, {
            'definition' => 'ssl_rsa_certificate'
          }, {
            'definition' => 'skip_cert_verify'
          }]
        }, {
          'type' => 'cloud_controller',
          'properties' => [{
            'definition' => 'system_domain'
          }, {
            'definition' => 'apps_domain'
          }, {
            'definition' => 'max_file_size'
          }]
        }]
      },
      product_settings_endpoint: '/components/cf/forms/ha_proxy/edit'
    }
  end

  def cf_1_3
    {
      name: 'cf',
      status: 'available',
      config: {
        'product_version' => '1.2.0.0',
        'type' => 'cf',
        'jobs' => [{
          'type' => 'ha_proxy',
          'properties' => [{
            'definition' => 'static_ips'
          }, {
            'definition' => 'ssl_rsa_certificate'
          }, {
            'definition' => 'skip_cert_verify'
          }]
        }, {
          'type' => 'cloud_controller',
          'properties' => [{
            'definition' => 'system_domain'
          }, {
            'definition' => 'apps_domain'
          }, {
            'definition' => 'max_file_size'
          }]
        }, {
          'type' => 'cloud_controller_worker',
          'resources' => [{
            'definition' => 'ephemeral_disk'
          }]
        }, {
          'type' => 'clock_global',
          'resources' => [{
            'definition' => 'ephemeral_disk'
          }]
        }],
        'properties' => [
          {
            'definition' => 'smtp_from'
          }, {
            'definition' => 'smtp_address'
          }, {
            'definition' => 'smtp_port'
          }, {
            'definition' => 'smtp_credentials'
          }, {
            'definition' => 'smtp_enable_starttls_auto'
          }, {
            'definition' => 'smtp_auth_mechanism'
          }
        ]
      },
      product_settings_endpoint: '/components/cf/forms/ha_proxy/edit'
    }
  end
end
# rubocop:enable Metrics/MethodLength
# rubocop:enable Metrics/ClassLength
