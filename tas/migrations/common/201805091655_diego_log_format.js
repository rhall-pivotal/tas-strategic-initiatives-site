exports.migrate = function(input) {
  var properties = input.properties;

  properties[".properties.diego_log_timestamp_format"] = { "value": "unix-epoch" };

  return input;
};
