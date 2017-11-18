exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.container_networking_log_traffic']['value'] == 'enable') {
    properties['.properties.container_networking_interface_plugin.silk.enable_log_traffic'] = { value: true };
    delete properties['.properties.container_networking_log_traffic']
  }

  return input;
};
