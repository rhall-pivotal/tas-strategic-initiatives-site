require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../201808151604_enforce_garden_shed_unselected.js")

describe("Ensure user has deployed without Garden Shed enabled", function() {
  it("raises an error if shed is enabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.enable_grootfs": { "value": false } } }
      )
    }).should.throw('attempt to upgrade to PCF Isolation Segment 2.3+ with GrootFS disabled, please enable GrootFS prior to upgrade by checking "Enable the GrootFS container image plugin for Garden RunC" in "Application Containers"');
  });

  it("no-ops if GrootFS is selected", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.enable_grootfs": { "value": true } } }
      )
    }).should.not.throw();
  });
});
