exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.enable_cf_metric_name'] = {
      value: true
  };

  return input;
}
