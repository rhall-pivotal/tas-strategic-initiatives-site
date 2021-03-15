exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.mysql_backups']['value'] == 'enable') {
    properties['.properties.mysql_backups.enable.backup_all_masters'] = {
        value: false
    };
  }

  return input;
};
