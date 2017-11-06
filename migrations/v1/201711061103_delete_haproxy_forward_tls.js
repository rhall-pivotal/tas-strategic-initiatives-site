exports.migrate = function(input) {
  if( input.properties['.properties.haproxy_forward_tls.enable.backend_ca'] ) {
    delete input.properties['.properties.haproxy_forward_tls.enable.backend_ca'];
  }
  return input;
};
