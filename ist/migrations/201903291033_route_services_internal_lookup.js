exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.route_services_internal_lookup'] = { value: true }

  return input;
};
