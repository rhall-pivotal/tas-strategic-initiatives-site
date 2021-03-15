exports.migrate = function(input) {
  var properties = input.properties;

  // Set selector until we get a fix for the OM bug
  properties['.properties.mysql_backups'] = {
      value: 'disable'
  };

  return input;
};
