exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.container_networking.disable.garden_network_pool'] = properties['.diego_cell.garden_network_pool'];

  return input;
};
