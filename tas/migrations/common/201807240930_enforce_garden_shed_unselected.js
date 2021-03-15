exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.enable_grootfs']['value'] === false) {
    abortMigration('attempt to upgrade to PAS 2.3+ with GrootFS disabled, please enable GrootFS prior to upgrade by checking "Enable the GrootFS container image plugin for Garden RunC" in "Application Containers"');
  }
  return input;
};
