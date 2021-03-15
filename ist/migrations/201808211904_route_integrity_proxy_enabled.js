exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.rep_proxy_enabled'] &&
    properties['.properties.rep_proxy_enabled']['value'] === true) {
    properties['.properties.route_integrity'] = { "value": "tls_verify" };
  }
  else {
    properties['.properties.route_integrity'] = { "value": "do_not_verify" };
  }

  delete properties['.properties.rep_proxy_enabled'];

  return input;
};
