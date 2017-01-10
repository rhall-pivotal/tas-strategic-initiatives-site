exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.mysql_monitor.recipient_email']) {
    if (properties['.mysql_monitor.recipient_email']['value'] == 'Fill in your desired email address') {
      properties['.mysql_monitor.recipient_email']['value'] = null;
    }
  }

  return input;
};
