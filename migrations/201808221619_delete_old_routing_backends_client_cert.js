exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.routing_backends_client_cert'];

  return input;
};
