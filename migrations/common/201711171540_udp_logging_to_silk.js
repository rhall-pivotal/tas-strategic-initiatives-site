exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.container_networking_log_traffic']['value'] == 'enabled') {
    properties['.properties.container_networking_interface_plugin.silk.iptables_accepted_udp_logs_per_sec'] = properties['.properties.container_networking_log_traffic.enable.iptables_accepted_udp_logs_per_sec']
    delete properties['.properties.container_networking_log_traffic.enable.iptables_accepted_udp_logs_per_sec']
  }

  return input;
};
