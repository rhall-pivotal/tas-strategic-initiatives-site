exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.container_networking_log_traffic'];

  return input;
};
