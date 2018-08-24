exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.garden_persistent_image_list_grootfs'];

  return input;
};
