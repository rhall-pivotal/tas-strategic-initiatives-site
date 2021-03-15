exports.migrate = function(input) {
  var properties = input.properties;

  properties['.diego_database.skip_consul_locks'] = { value: false };

  return input;
};
