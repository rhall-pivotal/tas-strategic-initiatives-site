require("./spec_helper.js");

const migration = require("../201803221053_delete_rep_server_cert.js");

describe("Remove pre-V2 rep agent certificate", function() {
  original_hash = {
    properties: {
      ".properties.rep_server_cert": {
        "type": "rsa_cert_credentials",
        "configurable": false,
        "credential": true,
        "value": {
          "private_key_pem": "some-private-key-pem",
          "cert_pem": "some-cert-pem"
        },
        "optional": false
      }
    }
  };

  migrated_hash = {
    properties: {},
  };

  it("removes rep_server_cert", function(){
    migration.migrate(original_hash).should.deep.equal(migrated_hash);
  });
});
