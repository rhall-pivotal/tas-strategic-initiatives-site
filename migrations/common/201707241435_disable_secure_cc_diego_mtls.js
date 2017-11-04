exports.migrate = function(input) {
  var properties = input.properties;

  properties['.cloud_controller.secure_diego_communication'] = { value: false };

  return input;
};
