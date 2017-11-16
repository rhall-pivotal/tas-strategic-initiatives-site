exports.migrate = function(input) {
  var properties = input.properties;

  if ( properties['.properties.credhub_key_encryption_password'] ) {
    properties['.properties.credhub_key_encryption_passwords'] = {
      value: [
        {
          guid: { value: generateGuid() },
          name: { value: 'Key' },
          key: properties['.properties.credhub_key_encryption_password'],
          primary: { value: true }
        }
      ]
    };

    delete properties['.properties.credhub_key_encryption_password'];
  }

  return input;
};
