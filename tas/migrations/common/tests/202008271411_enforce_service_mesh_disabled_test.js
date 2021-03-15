require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../202008271411_enforce_service_mesh_disabled.js");

describe("Ensure user has deployed without Service Mesh enabled", function() {

  it("raises an error if service mesh is enabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.istio": { "value": "enable" } } }
      )
    }).should.throw('attempt to upgrade to PAS 2.11+ with Service Mesh enabled, please disable Service Mesh prior to upgrade by setting ".properties.istio" to "disable"');
  });

  it("no-ops if Service Mesh is disabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.istio": { "value": "disable" } } }
      )
    }).should.not.throw();
  });
});
