require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../202008312231_remove_dynamic_egress.js");

describe("Ensure user has deployed without Dynamic Egress enabled", function() {

  it("raises an error if Dynamic Egress is enabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.experimental_dynamic_egress_enforcement": { "value": true } } }
      )
    }).should.throw('attempt to upgrade to TAS 2.11+ with Dynamic Egress enabled, please disable Dynamic Egress to upgrade. You can do this by setting `.properties.experimental_dynamic_egress_enforcement` to `false` using the OM CLI.');
  });

  it("no-ops if Dynamic Egress is Disabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.experimental_dynamic_egress_enforcement": { "value": false } } }
      )
    }).should.not.throw();
  });
});
