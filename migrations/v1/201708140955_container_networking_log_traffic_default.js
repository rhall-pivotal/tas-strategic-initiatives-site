exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.container_networking_log_traffic'] = {
    value: 'disable'
  };

  return input;
};
