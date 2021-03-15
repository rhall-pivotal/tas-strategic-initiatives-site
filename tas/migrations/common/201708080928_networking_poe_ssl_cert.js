exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.networking_point_of_entry'] ) {
    if( properties['.properties.networking_point_of_entry']['value'] == 'haproxy' ) {
      properties['.properties.networking_poe_ssl_cert'] = properties['.properties.networking_point_of_entry.haproxy.ssl_rsa_certificate'];
    } else if( properties['.properties.networking_point_of_entry']['value'] == 'external_ssl' ) {
      properties['.properties.networking_poe_ssl_cert'] = properties['.properties.networking_point_of_entry.external_ssl.ssl_rsa_certificate'];
    }
  }

  return input;
};
