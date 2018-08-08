exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.uaa.ldap.server_ssl_cert_alias'];

  return input;
};
