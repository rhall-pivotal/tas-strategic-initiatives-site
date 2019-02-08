
require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../201902080956_internal_domain_to_collection.js");

describe("changing internal_domain to collection of internal_domains", function() {
    it("migrates the single string property to a collection", function() {
      migration.migrate(
        { properties: {".properties.cf_networking_internal_domain": {"value": "apps.meow.internal" } } }
      ).should.deep.equal(
        { properties: { ".properties.cf_networking_internal_domains": [{ "name": "apps.meow.internal" }] } }
      );
  });
});
