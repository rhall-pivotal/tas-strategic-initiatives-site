exports.migrate = function(input) {
  var properties = input.properties;
  delete properties['.properties.secure_diego_communication']
  return input;
};
