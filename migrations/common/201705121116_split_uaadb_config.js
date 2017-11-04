exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.uaa_database'] = {
      value: properties['.properties.system_database']['value']
  };

  if (properties['.properties.system_database']['value'] == 'external') {
    properties['.properties.uaa_database.external.uaa_username'] = properties['.properties.system_database.external.uaa_username']
    properties['.properties.uaa_database.external.uaa_password'] = properties['.properties.system_database.external.uaa_password']
    properties['.properties.uaa_database.external.host'] = properties['.properties.system_database.external.host']
    properties['.properties.uaa_database.external.port'] = properties['.properties.system_database.external.port']
  }

  return input;
};
