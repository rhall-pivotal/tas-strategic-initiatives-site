exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.uaa.saml.entity_id_override'] ) {
    properties['.properties.saml_entity_id_override'] = properties['.properties.uaa.saml.entity_id_override'];
  }

  return input;
};
