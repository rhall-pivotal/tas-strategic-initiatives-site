exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.autoscale_metric_bucket_count'] = { "value": 120 };

  return input;
};
