exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.uaa.saml.signature_algorithm'] ) {
    if( properties['.properties.uaa.saml.signature_algorithm']['value'] == "SHA1" ) {
      properties['.properties.saml_signature_algorithm'] = {
        value: "SHA256"
      };
    }
  }

  return input;
};
