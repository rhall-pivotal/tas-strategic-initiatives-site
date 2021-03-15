exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.router_keepalive_connections'] = {
    value: 'enable'
  };

  return input;
};
