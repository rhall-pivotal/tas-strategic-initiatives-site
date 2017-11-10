exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.uaa.saml.signature_algorithm'] ) {
    properties['.properties.saml_signature_algorithm'] = properties['.properties.uaa.saml.signature_algorithm'];
  }

  return input;
};
