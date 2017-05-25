exports.migrate = function(input) {
    var properties = input.properties;
    var oldValue = properties['.uaa.service_provider_key_password']['value']

    properties['.uaa.service_provider_key_password'] = {
      value: {
        secret: oldValue
      }
    };

    return input;
};
