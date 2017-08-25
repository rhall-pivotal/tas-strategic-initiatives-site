exports.migrate = function(input) {
  var properties = input.properties;

  var key = properties['.notifications.encryption_credentials']['value']['password'];
  properties['.notifications.encryption_key'] = {
    value: {
      secret: key
    }
  };

  delete properties['.notifications.encryption_credentials']

  return input;
};
