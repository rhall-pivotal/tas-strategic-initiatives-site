exports.migrate = function(input) {
  var properties = input.properties;

  properties['.autoscaling.encryption_key']['value'] = properties['.autoscaling.encryption_key']['password'];

  return input;
};
