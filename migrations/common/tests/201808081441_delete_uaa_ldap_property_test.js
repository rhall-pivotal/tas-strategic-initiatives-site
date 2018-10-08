require("tap").mochaGlobals()
const should = require("chai").should()
const migration = require("../201808081441_delete_uaa_ldap_property.js")

describe("Remove UAA LDAP Server SSL Cert Alias", function() {
  original_hash = {
    properties: {
      ".properties.uaa.ldap.server_ssl_cert_alias": {
        "type": "string",
        "configurable": true,
        "optional": true
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes uaa.ldap.server_ssl_cert_alias", function(){
    migration.migrate(original_hash).should.deep.equal(migrated_hash);
  });
});
