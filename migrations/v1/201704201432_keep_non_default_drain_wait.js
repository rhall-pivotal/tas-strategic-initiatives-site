exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.router.drain_wait'] != 0) {
    properties['.router.drain_wait'] = properties['.router.drain_wait']
  }

  return input;
};
