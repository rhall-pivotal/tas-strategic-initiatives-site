exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.networking_point_of_entry']['value'] == 'terminate_at_router' ) {
    properties['.properties.networking_poe_ssl_cert'] = properties['.properties.networking_point_of_entry.terminate_at_router.ssl_rsa_certificate'];
    delete properties['.properties.networking_point_of_entry.terminate_at_router.ssl_rsa_certificate'];
  }

  return input;
};
