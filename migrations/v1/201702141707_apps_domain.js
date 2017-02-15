exports.migrate = function(input) {
  var properties = input.properties;

  properties['.cloud_controller.primary_apps_domain'] = properties['.cloud_controller.apps_domain']

  return input;
};
