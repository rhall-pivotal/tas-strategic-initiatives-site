exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.uaa.service_provider_key_password']) {
    var oldValue = properties['.uaa.service_provider_key_password']['value']

    properties['.uaa.service_provider_key_password'] = {
      value: {
        secret: oldValue
      }
    }
  };

  return input;
};
