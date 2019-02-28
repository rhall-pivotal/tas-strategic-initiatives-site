require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../201902280953_prevent_upgrade_on_cf_metrics.js");

describe("Ensure deprecated 'cf' metrics name is not selected", function() {
  context("when enable_cf_metric_name is selected", function() {
    it("raises an error", function() {
      (function() {
        migration.migrate(
          { properties: { ".properties.enable_cf_metric_name": {"value": true } } }
        )
      }).should.throw("Attempt to upgrade to PAS for Windows 2.5+ with the deprecated 'cf' metrics name selected. Please uncheck the 'Use \"cf\" as deployment name in emitted metrics' option in the Advanced Features tab before attempting to upgrade.");
    });
  });

  context("when enable_cf_metric_name is not selected", function() {
    it("does nothing", function() {
      (function() {
        migration.migrate(
          { properties: { ".properties.enable_cf_metric_name": {"value": false } } }
        )
      }).should.not.throw();
    });
  });
});
