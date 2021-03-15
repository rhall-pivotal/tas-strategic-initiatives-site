exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.enable_service_discovery_for_apps']['value'] === false) {
    properties['.properties.cf_networking_internal_domain'] = { "value": "" };
  }

  delete properties['.properties.enable_service_discovery_for_apps'];

  return input;
};
