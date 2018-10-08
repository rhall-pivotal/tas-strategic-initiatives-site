require("tap").mochaGlobals()
const should = require("chai").should()

const migration = require("../201808231922_delete_garden_persistent_image_list_grootfs.js")

describe("Remove garden_persistent_image_list_grootfs property", function() {
  original_hash = {
    properties: {
      ".properties.garden_persistent_image_list_grootfs": {
        "value": ["/var/vcap/packages/cflinuxfs2/rootfs.tar"]
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes garden_persistent_image_list_grootfs", function() {
    migration.migrate(original_hash).should.deep.equal(migrated_hash);
  });
});
