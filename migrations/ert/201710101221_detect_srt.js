exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Small Footprint Elastic Runtime cannot be upgraded to Elastic Runtime";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "SRT" ) {
      abortMigration(errorMessage);
    }
  }

  return input;
};
