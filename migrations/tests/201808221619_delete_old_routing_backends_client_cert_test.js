require("tap").mochaGlobals()
const should = require("chai").should()

const migration = require("../201808221619_delete_old_routing_backends_client_cert")

describe("Remove old GoRouter Backend client cert", function() {
  original_hash = {
    properties: {
      ".properties.routing_backends_client_cert": {
        "type": "rsa_cert_credentials"
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes routing_backends_client_cert", function() {
    migration.migrate(original_hash).should.deep.equal(migrated_hash);
  });
});
