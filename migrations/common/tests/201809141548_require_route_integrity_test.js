require("tap").mochaGlobals()
const should = require("should")

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../201809141548_require_route_integrity.js")

describe("Ensure route integrity is enabled", function() {

  context("when route integrity is disabled", function() {
    it("raises an error", function(){
      (function () {
        migration.migrate(
          { properties: { ".properties.route_integrity": { "value": "do_not_verify" } } }
        )
      }).should.throw('attempt to upgrade to PAS 2.4+ with route integrity disabled, please enable route integrity prior to upgrade under "Router application identity verification" in "Application Containers"');
    });
  });

  context("when TLS route integrity is enabled", function() {
    it("does nothing", function(){
      (function () {
        migration.migrate(
          { properties: { ".properties.route_integrity": { "value": "tls_verify" } } }
        )
      }).should.not.throw();
    });
  });

  context("when mutual TLS route integrity is enabled", function() {
    it("does nothing", function(){
      (function () {
        migration.migrate(
          { properties: { ".properties.route_integrity": { "value": "mutual_tls_verify" } } }
        )
      }).should.not.throw();
    });
  });

});
