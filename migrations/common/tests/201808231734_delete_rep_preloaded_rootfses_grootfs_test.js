require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201808231734_delete_rep_preloaded_rootfses_grootfs.js")

describe("Remove rep_preloaded_rootfses_grootfs property", function() {
  original_hash = {
    properties: {
      ".properties.rep_preloaded_rootfses_grootfs": ["cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs.tar"],
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes routing_backends_client_cert", function() {
    migration.migrate(original_hash).should.deepEqual(migrated_hash);
  });
});
