require("tap").mochaGlobals()
const should = require("chai").should()
const migration = require("../201808211904_route_integrity_proxy_enabled.js")

describe("Route Integrity", function() {
  context("when the proxy enabled property is not defined", function() {
    it("sets route_integrity to do_not_verify", function(){
      migration.migrate(
        { properties: {} }
      ).should.deep.equal(
        { properties: { ".properties.route_integrity": { "value": "do_not_verify" } } }
      );
    });
  });
  context("when the proxy is disabled", function() {
    it("sets route_integrity to do_not_verify", function(){
      migration.migrate(
        { properties: { ".properties.rep_proxy_enabled": { "value": false } } }
      ).should.deep.equal(
        { properties: { ".properties.route_integrity": { "value": "do_not_verify" } } }
      );
    });
  });
  context("when the proxy is enabled", function() {
    it("sets route_integrity to tls_verify", function(){
      migration.migrate(
        { properties: { ".properties.rep_proxy_enabled": { "value": true } } }
      ).should.deep.equal(
        { properties: { ".properties.route_integrity": { "value": "tls_verify" } } }
      );
    });
  });
});
