exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.system_logging'] = {value: 'disabled'};

  if(properties['.properties.syslog_host']['value']) {
    properties['.properties.system_logging'] = {value: 'enabled'};
    properties['.properties.system_logging.enabled.host'] = properties['.properties.syslog_host'];
    properties['.properties.system_logging.enabled.port'] = properties['.properties.syslog_port'];
    properties['.properties.system_logging.enabled.protocol'] = properties['.properties.syslog_protocol'];
  }

  return input;
};
