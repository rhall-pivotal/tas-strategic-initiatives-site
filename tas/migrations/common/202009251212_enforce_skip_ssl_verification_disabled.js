exports.migrate = function(input) {
  var properties = input.properties;
  var errmsg = 'attempt to upgrade to PAS 2.11+ with Skip SSL Verification enabled, please disable Skip SSL Verification prior to upgrade by un-checking "Disable SSL certificate verification for this environment" under "Networking"'

  if (properties['.ha_proxy.skip_cert_verify']['value'] === true) {
    abortMigration(errmsg);
  }
  return input;
};
