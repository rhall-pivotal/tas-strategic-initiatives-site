exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.networking_point_of_entry'] ) {
    if( properties['.properties.networking_point_of_entry']['value'] == 'haproxy' ) {
      properties['.properties.routing_disable_http'] = properties['.properties.networking_point_of_entry.haproxy.disable_http'];
    }
  }

  return input;
};
