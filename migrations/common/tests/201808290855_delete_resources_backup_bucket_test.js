require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201808290855_delete_resources_backup_bucket.js")

describe("Remove resources_backup_bucket property", function() {
  original_hash = {
    properties: {
      ".properties.resources_backup_bucket": "some-resource-backup-bucket",
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes resources_backup_bucket", function() {
    migration.migrate(original_hash).should.deepEqual(migrated_hash);
  });
});
