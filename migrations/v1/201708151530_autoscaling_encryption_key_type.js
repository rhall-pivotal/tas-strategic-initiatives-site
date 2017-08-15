exports.migrate = function(input) {
  var properties = input.properties;

  properties['.autoscaling.encryption_key'] = {
    value: properties['.autoscaling.encryption_key.haproxy.max_buffer_size']['password']
  };

  return input;
};
