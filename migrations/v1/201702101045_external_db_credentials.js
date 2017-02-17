exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.system_database']['value'] == 'external') {
    if (properties['.properties.system_database.external.username']['value'] != '') {
      properties['.properties.system_database.external.account_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.account_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.app_usage_service_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.app_usage_service_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.autoscale_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.autoscale_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.ccdb_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.ccdb_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.diego_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.diego_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.networkpolicyserver_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.networkpolicyserver_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.nfsvolume_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.nfsvolume_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.notifications_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.notifications_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.routing_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.routing_password'] = properties['.properties.system_database.external.password']
      properties['.properties.system_database.external.uaa_username'] = properties['.properties.system_database.external.username']
      properties['.properties.system_database.external.uaa_password'] = properties['.properties.system_database.external.password']
    }
  }

  return input;
};
