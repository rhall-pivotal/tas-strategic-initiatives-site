exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.container_networking_vtep_port']) {
    properties['.properties.container_networking_interface_plugin.silk.vtep_port'] = properties['.properties.container_networking_vtep_port']
    delete properties['.properties.container_networking_vtep_port']
  }

  return input;
};
