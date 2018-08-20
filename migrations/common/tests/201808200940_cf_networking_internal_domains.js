require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201808200931_cf_networking_internal_domains.js")

describe("Set internal domains", function() {
  it("sets cf_networking_internal_domains to empty string on upgrade", function(){
    migration.migrate(
      { properties: { ".properties.enable_service_discovery_for_apps": { "value": false } } }
    ).should.deepEqual(
      { properties: { ".properties.cf_networking_internal_domain": { "value": "" } } }
    );
  });
});

describe("Remove enable_service_discovery_for_apps", function(){
  original_hash = {
    properties: {
      ".properties.enable_service_discovery_for_apps": {
        "type": "boolean",
        "configurable": true,
        "default": true
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes enable_service_discovery_for_apps", function(){
    migration.migrate(original_hash).should.deepEqual(migrated_hash);
  });
});
