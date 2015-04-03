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
    when '1.3'
      components << cf_1_3
    when '1.4'
      components << cf_1_4
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
          'instances' => [{
            'identifier' => 'instances'
          }],
          'resources' => [{
            'definition' => 'ephemeral_disk'
          },]
        }, {
          'type' => 'clock_global',
          'instances' => [{
            'identifier' => 'instances'
          }],
          'resources' => [{
            'definition' => 'ephemeral_disk'
          },]
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

  def cf_1_4
    {
      name: 'cf',
      status: 'available',
      config: {
        'product_version' => '1.4.0.0',
        'identifier' => 'cf',
        'jobs' => [{
          'identifier' => 'ha_proxy',
          'properties' => [{
            'identifier' => 'static_ips'
          }, {
            'identifier' => 'ssl_rsa_certificate'
          }, {
            'identifier' => 'skip_cert_verify'
          }]
        }, {
          'identifier' => 'cloud_controller',
          'properties' => [{
            'identifier' => 'system_domain'
          }, {
            'identifier' => 'apps_domain'
          }, {
            'identifier' => 'max_file_size'
          }]
        }, {
          'identifier' => 'cloud_controller_worker',
          'instances' => [{
            'identifier' => 'instances'
          }],
          'resources' => [{
            'identifier' => 'ephemeral_disk'
          },]
        }, {
          'identifier' => 'clock_global',
          'instances' => [{
            'identifier' => 'instances'
          }],
          'resources' => [{
            'identifier' => 'ephemeral_disk'
          },]
        }, {
          'identifier' => 'router'
        }],
        'properties' => [
          {
            'identifier' => 'smtp_from'
          }, {
            'identifier' => 'smtp_address'
          }, {
            'identifier' => 'smtp_port'
          }, {
            'identifier' => 'smtp_credentials'
          }, {
            'identifier' => 'smtp_enable_starttls_auto'
          }, {
            'identifier' => 'smtp_auth_mechanism'
          }, {
            'identifier' => 'logger_endpoint_port'
          }
        ]
      },
      product_settings_endpoint: '/components/cf/forms/ha_proxy/edit'
    }
  end
end
# rubocop:enable Metrics/MethodLength
# rubocop:enable Metrics/ClassLength
