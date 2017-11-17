exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.container_networking_network_cidr']) {
    properties['.properties.container_networking_interface_plugin.silk.network_cidr'] = properties['.properties.container_networking_network_cidr']
    delete properties['.properties.container_networking_network_cidr']
  }

  return input;
};
