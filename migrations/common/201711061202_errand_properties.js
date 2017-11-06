exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.push_usage_service_secret_token'] = properties['.push-usage-service.secret_token'];

  return input;
};
