exports.migrate = function(input) {
  var properties = input.properties;

  if( properties['.properties.router_forward_client_cert'] ) {
    if( properties['.properties.router_forward_client_cert']['value'] == 'always_forward' ) {
      properties['.properties.routing_tls_termination'] = {
        value: 'load_balancer'
      };
    } else if( properties['.properties.router_forward_client_cert']['value'] == 'forward' ) {
      properties['.properties.routing_tls_termination'] = {
        value: 'ha_proxy'
      };
    } else if( properties['.properties.router_forward_client_cert']['value'] == 'sanitize_set' ) {
      properties['.properties.routing_tls_termination'] = {
        value: 'router'
      };
    }
  }

  return input;
};
