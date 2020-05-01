exports.migrate = function(input) {
  const properties = input.properties;

  if( properties['.properties.diego_log_timestamp_format']['value'] == 'unix-epoch' ) {
    properties['.properties.logging_timestamp_format'] = { value: 'deprecated' }
  }

  if( properties['.properties.diego_log_timestamp_format']['value'] == 'rfc3339' ) {
    properties['.properties.logging_timestamp_format'] = { value: 'rfc3339' }
  }

  delete properties['.properties.diego_log_timestamp_format'];

  return input;
};
