exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.route_services']['value'] == 'enable' ) {
    properties['.properties.route_services.enable.internal_lookup'] = { value: true }
  }

  return input;
};
