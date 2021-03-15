exports.migrate = function(input) {
  var properties = input.properties;

  // Set selector until we get a fix for the OM bug
  properties['.properties.garden_disk_cleanup'] = {
      value: 'routine'
  };

  return input;
};
