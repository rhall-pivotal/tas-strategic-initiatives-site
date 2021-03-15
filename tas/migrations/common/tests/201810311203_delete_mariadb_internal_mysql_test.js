require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../201810311203_delete_mariadb_internal_mysql.js");

describe("Ensure Internal MySQL MariaDB Cluster is not selected", function() {
  context("when internal_mysql is selected", function() {
    it("raises an error", function() {
      (function() {
        migration.migrate(
          { properties: { ".properties.system_database": {"value": "internal_mysql" } } }
        )
      }).should.throw("Attempt to upgrade to PAS 2.4+ with the deprecated Internal MySQL MariaDB Cluster selected. Please consult the documentation for guidance on how to migrate to the Internal MySQL Percona XtraDB Cluster (PXC) database prior to upgrade.");
    });
  });

  context("when internal_pxc is selected", function() {
    it("does nothing", function() {
      (function() {
        migration.migrate(
          { properties: { ".properties.system_database": {"value": "internal_pxc" } } }
        )
      }).should.not.throw();
    });
  });

  context("when external is selected", function() {
    it("does nothing", function() {
      (function() {
        migration.migrate(
          { properties: { ".properties.system_database": {"value": "external" } } }
        )
      }).should.not.throw();
    });
  });
});
