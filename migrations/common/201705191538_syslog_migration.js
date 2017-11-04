exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.syslog_tls'] = {value: 'disabled'};

  return input;
};
