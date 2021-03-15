exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.push_apps_manager_square_logo'] ) {
    properties['.properties.push_apps_manager_favicon'] = properties['.properties.push_apps_manager_square_logo'];
  }

  return input;
};
