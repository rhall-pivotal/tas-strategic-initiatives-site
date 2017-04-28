exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['encrypt_key']['value'] ) {
    properties['.db_encryption_credentials.password'] = properties['encrypt_key']['value']
  }

  return input
};
