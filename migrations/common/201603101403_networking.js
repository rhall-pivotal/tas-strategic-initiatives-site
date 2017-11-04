exports.migrate = function(input) {
    var properties = input.properties;

    if (properties['.ha_proxy.static_ips']['value']) {
        properties['.properties.networking_point_of_entry'] = {
            value: 'haproxy'
        };

        properties['.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate'] = properties['.ha_proxy.ssl_rsa_certificate'];
        properties['.properties.networking_point_of_entry.haproxy.disable_http'] = properties['.ha_proxy.disable_http'];
        properties['.properties.networking_point_of_entry.haproxy.ssl_ciphers'] = properties['.ha_proxy.ssl_ciphers'];
    } else if (properties['.router.enable_ssl']['value']) {
        properties['.properties.networking_point_of_entry'] = {
            value: 'external_ssl'
        };

        properties['.properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate'] = properties['.ha_proxy.ssl_rsa_certificate'];
        properties['.properties.networking_point_of_entry.external_ssl.ssl_ciphers'] = properties['.router.ssl_ciphers'];
    } else {
        properties['.properties.networking_point_of_entry'] = {
            value: 'external_non_ssl'
        };
    }

    // Always enable route services by default.
    properties['.properties.route_services'] = {
        value: 'enable'
    };

    if (properties['.ha_proxy.disable_http']['value']) {
        properties['.properties.route_services.enable.ignore_ssl_cert_verification'] = properties['.ha_proxy.disable_http'];
    }

    return input;
}
