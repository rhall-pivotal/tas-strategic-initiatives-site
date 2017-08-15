exports.migrate = function(input) {
  var properties = input.properties;

  if ( properties['.properties.container_networking'] ) {
    if ( properties['.properties.container_networking']['value'] == 'enable' ) {
      properties['.properties.container_networking_vtep_port'] = properties['.properties.container_networking.enable.vtep_port'];
      delete properties['.properties.container_networking.enable.vtep_port'];
    }
  }

  return input;
};
