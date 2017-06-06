exports.migrate = function(input) {
  var properties = input.properties;

  // Disables c2c on upgrades from 1.10
  properties['.properties.container_networking'] = {
      value: 'disable'
  };

  return input;
};
