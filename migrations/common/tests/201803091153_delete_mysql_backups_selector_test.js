require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201803091153_delete_mysql_backups_selector.js")

describe("Remove Automated Backup Configuration", function() {
  original_hash = {
    properties: {
      ".properties.mysql_backups": { "option1": "some-value1", "option2": "some-value2" },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes mysql_backup", function(){
    migration.migrate(original_hash).should.deepEqual(migrated_hash);
  });

});
