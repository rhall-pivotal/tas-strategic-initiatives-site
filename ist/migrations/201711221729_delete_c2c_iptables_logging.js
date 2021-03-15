exports.migrate = function(input) {
  delete input.properties['.properties.c2c_iptables_logging.enabled.denied_logging_interval'];
  delete input.properties['.properties.c2c_iptables_logging.enabled.udp_logging_interval'];
  delete input.properties['.properties.c2c_iptables_logging'];
  return input;
};
