require("./spec_helper.js");

const migration = require("../201902080956_internal_domain_to_collection.js");

describe("changing internal_domain to collection of internal_domains", function() {
    it("migrates the single string property to a collection", function() {
      var migratedDomains = migration.migrate(
        { properties: {".properties.cf_networking_internal_domain": {"value": "apps.meow.internal" } } }
      )['properties']['.properties.cf_networking_internal_domains']['value'];

      migratedDomains.length.should.equal(1);
      migratedDomains[0]['name']['value'].should.equal('apps.meow.internal');
  });
});
