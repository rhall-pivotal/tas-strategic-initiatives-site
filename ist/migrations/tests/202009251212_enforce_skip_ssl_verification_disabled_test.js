require("./spec_helper.js");

abortMigration = function(msg) {
  throw new Error(msg);
};

const migration = require("../202009251212_enforce_skip_ssl_verification_disabled.js");

describe("Ensure user has deployed without 'Skip SSL Verification' enabled", function() {

  it("raises an error if Skip SSl Verification is enabled", function(){

    (function () {
      migration.migrate(
        { properties: { ".properties.skip_cert_verify": { "value": true } } }
      )
    }).should.throw('attempt to upgrade to IST 2.11+ with Skip SSL Verification enabled, please disable Skip SSL Verification prior to upgrade by un-checking "Disable SSL certificate verification for this environment" under "Networking"');
  });

  it("no-ops if Skip SSL Verification is disabled", function(){
    (function () {
      migration.migrate(
        { properties: { ".properties.skip_cert_verify": { "value": false } } }
      )
    }).should.not.throw();
  });
});
