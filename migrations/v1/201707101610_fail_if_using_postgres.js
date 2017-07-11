exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.system_database']['value'] == 'internal') {
    abortMigration("This release does not support an Internal Postgres database. Please reconfigure your system databases before proceeding");
  }

  if (properties['.properties.uaa_database']['value'] == 'internal') {
    abortMigration("This release does not support an Internal Postgres database. Please reconfigure your UAA database before proceeding");
  }

  return input;
};
