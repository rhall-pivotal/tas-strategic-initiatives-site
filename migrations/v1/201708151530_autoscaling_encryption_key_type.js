exports.migrate = function(input) {
  var properties = input.properties;

  var key = properties['.autoscaling.encryption_key']['password'];
  properties['.autoscaling.encryption_key'] = {
    value: {
      secret: key
    }
  };

  return input;
};
