exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "VMware Tanzu Application Service cannot be upgraded to Small Footprint VMware Tanzu Application Service";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "ERT" ) {
      abortMigration(errorMessage);
    }
  } else {
    abortMigration(errorMessage);
  }

  return input;
};
