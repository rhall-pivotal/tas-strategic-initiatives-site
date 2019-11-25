exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Pivotal Application Service cannot be upgraded to Pivotal Platform Small Footprint";

  if( properties['.properties.tile_name'] ) {
    if( properties['.properties.tile_name']['value'] == "ERT" ) {
      abortMigration(errorMessage);
    }
  } else {
    abortMigration(errorMessage);
  }

  return input;
};
