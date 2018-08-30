exports.migrate = function(input) {
  var properties = input.properties;
  delete properties['.properties.resources_backup_bucket']
  return input;
};
