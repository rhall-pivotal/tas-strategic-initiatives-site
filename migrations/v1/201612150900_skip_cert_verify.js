exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.skip_cert_verify'] = properties['.ha_proxy.skip_cert_verify'];

  return input;
}
