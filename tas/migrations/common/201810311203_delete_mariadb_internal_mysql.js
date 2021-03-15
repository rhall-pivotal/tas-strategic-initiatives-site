exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Attempt to upgrade to PAS 2.4+ with the deprecated Internal MySQL MariaDB Cluster selected. Please consult the documentation for guidance on how to migrate to the Internal MySQL Percona XtraDB Cluster (PXC) database prior to upgrade.";

  if(properties['.properties.system_database']['value'] == "internal_mysql") {
      abortMigration(errorMessage);
  }

  return input;
}
