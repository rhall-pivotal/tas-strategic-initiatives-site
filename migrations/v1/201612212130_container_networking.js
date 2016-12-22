exports.migrate = function(input) {
  var properties = input.properties;

  // Set selector until we get a fix for the OM bug
  properties['.properties.container_networking'] = {
      value: 'disable'
  };

  return input;
};
