require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201808231734_delete_rep_preloaded_rootfses_grootfs.js")

describe("Remove rep_preloaded_rootfses_grootfs property", function() {
  original_hash = {
    properties: {
      ".properties.rep_preloaded_rootfses_grootfs": {
        "value": ["cflinuxfs2:/var/vcap/packages/cflinuxfs2/rootfs.tar"]
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes rep_preloaded_rootfses_grootfs", function() {
    migration.migrate(original_hash).should.deepEqual(migrated_hash);
  });
});
