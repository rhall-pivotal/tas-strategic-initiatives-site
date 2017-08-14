exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.networking_point_of_entry'] ) {
    if( properties['.properties.networking_point_of_entry']['value'] == 'haproxy' ) {
      properties['.properties.haproxy_max_buffer_size'] = properties['.properties.networking_point_of_entry.haproxy.max_buffer_size'];
    }
  }

  return input;
};
