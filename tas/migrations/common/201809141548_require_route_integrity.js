exports.migrate = function(input) {
  var properties = input.properties;

  if (properties['.properties.route_integrity']) {
    if (properties['.properties.route_integrity']['value'] === "do_not_verify") {
      abortMigration('attempt to upgrade to PAS 2.4+ with route integrity disabled, please enable route integrity prior to upgrade under "Router application identity verification" in "Application Containers"');
    }
  }
  return input;
};
