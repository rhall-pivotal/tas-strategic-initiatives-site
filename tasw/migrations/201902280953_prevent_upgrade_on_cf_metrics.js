exports.migrate = function(input) {
  var properties = input.properties;
  var errorMessage  = "Attempt to upgrade to PAS for Windows 2.5+ with the deprecated 'cf' metrics name selected. Please uncheck the 'Use \"cf\" as deployment name in emitted metrics' option in the Advanced Features tab before attempting to upgrade.";

  if(properties['.properties.enable_cf_metric_name']['value'] == true) {
      abortMigration(errorMessage);
  }

  return input;
}
