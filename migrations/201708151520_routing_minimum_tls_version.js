exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.routing_minimum_tls_version'] = {
    value: 'tls_v1_2'
  };

  return input;
};
