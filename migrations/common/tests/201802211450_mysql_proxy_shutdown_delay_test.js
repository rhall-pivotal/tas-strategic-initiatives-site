require("tap").mochaGlobals()
const should = require("chai").should()
const migration = require("../201802211450_mysql_proxy_shutdown_delay.js")

describe("MySQL proxy shutdown delay", function() {

  context("when the existing shutdown delay is zero", function() {
    it("increases the shutdown delay to 30 seconds", function(){
      migration.migrate(
        { properties: { ".mysql_proxy.shutdown_delay": { "value": 0 } } }
      ).should.deep.equal(
        { properties: { ".mysql_proxy.shutdown_delay": { "value": 30 } } }
      );
    });
  });

  context("when the existing shutdown delay is non-zero", function() {
    it("retains the existing value", function(){
      migration.migrate(
        { properties: { ".mysql_proxy.shutdown_delay": { "value": 60 } } }
      ).should.deep.equal(
        { properties: { ".mysql_proxy.shutdown_delay": { "value": 60 } } }
      );
    });
  });

});
