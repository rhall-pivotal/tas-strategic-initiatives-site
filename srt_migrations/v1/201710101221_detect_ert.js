exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Elastic Runtime cannot be upgraded to Small Footprint Elastic Runtime";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "ERT" ) {
      abortMigration(errorMessage);
    }
  } else {
    abortMigration(errorMessage);
  }

  return input;
};
