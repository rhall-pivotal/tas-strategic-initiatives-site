exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.mysql_backups'];

  return input;
};
