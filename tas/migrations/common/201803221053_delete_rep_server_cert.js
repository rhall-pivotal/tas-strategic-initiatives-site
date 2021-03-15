exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.rep_server_cert'];

  return input;
};
