exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "PCF Small Footprint cannot be upgraded to Pivotal Application Service";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "SRT" ) {
      abortMigration(errorMessage);
    }
  }

  return input;
};
