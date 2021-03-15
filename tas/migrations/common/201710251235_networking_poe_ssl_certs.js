exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.networking_poe_ssl_certs'] = {
    value: [
      {
        guid: { value: generateGuid() },
        name: { value: 'Certificate' },
        certificate: properties['.properties.networking_poe_ssl_cert']
      }
    ]
  };

  delete properties['.properties.networking_poe_ssl_cert'];

  return input;
};
