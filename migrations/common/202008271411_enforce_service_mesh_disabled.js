exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.istio']['value'] == 'enable' ) {
    abortMigration('attempt to upgrade to PAS 2.11+ with Service Mesh enabled, please disable Service Mesh prior to upgrade by setting ".properties.istio" to "disable"');
  }
  return input;
};
