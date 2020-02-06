exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Small Footprint VMware Tanzu Application Service cannot be upgraded to VMware Tanzu Application Service";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "SRT" ) {
      abortMigration(errorMessage);
    }
  }

  return input;
};
