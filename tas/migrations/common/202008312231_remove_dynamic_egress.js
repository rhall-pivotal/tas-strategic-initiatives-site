exports.migrate = function(input) {
var properties = input.properties;

  if (properties['.properties.experimental_dynamic_egress_enforcement']['value'] === true) {
    abortMigration('attempt to upgrade to TAS 2.11+ with Dynamic Egress enabled, please disable Dynamic Egress to upgrade. You can do this by setting `.properties.experimental_dynamic_egress_enforcement` to `false` using the OM CLI.');
  }
  return input;
};
