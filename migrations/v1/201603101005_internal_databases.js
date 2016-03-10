exports.migrate = function(input) {
  var properties = input['.properties'];

  // Migrate mixed internal databases to MySQL-only
  if (properties['.properties.system_database'].value == "internal") {
    properties['.properties.system_database'].value = "internal_mysql";
  }

  return input;
};
