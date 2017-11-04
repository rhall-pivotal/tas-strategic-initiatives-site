exports.migrate = function(input) {
  var properties = input.properties;

  properties['.push-usage-service.secret_token'] = properties['.push-apps-manager.secret_token'];

  return input;
};
