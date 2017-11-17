exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.diego_cell.garden_network_mtu'] ) {
    if( properties['.properties.container_networking_interface_plugin']['value'] == 'silk' ) {
      properties['.properties.container_networking_interface_plugin.silk.network_mtu'] = properties['.diego_cell.garden_network_mtu'];
    }

    delete properties['.diego_cell.garden_network_mtu'];
  }

  return input;
};
