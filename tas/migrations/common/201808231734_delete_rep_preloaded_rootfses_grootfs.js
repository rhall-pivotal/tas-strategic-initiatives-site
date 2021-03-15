exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.rep_preloaded_rootfses_grootfs'];

  return input;
};
