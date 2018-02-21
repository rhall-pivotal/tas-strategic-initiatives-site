exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.mysql_proxy.shutdown_delay'] ) {
    if( properties['.mysql_proxy.shutdown_delay']['value'] == 0 ) {
      properties['.mysql_proxy.shutdown_delay']['value'] = 30;
    }
  }

  return input;
};
