exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.container_networking_interface_plugin'] = {
    value: 'silk'
  };

  return input;
};
