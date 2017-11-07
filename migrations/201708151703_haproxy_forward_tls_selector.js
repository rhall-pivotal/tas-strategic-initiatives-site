exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.haproxy_forward_tls'] = {
    value: 'enable'
  };

  return input;
};
