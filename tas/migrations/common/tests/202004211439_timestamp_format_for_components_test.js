require("./spec_helper.js");

const migration = require("../202004211439_timestamp_format_for_components.js");

describe("Remove Diego Timestamp Property", function() {
  original_hash_epoch = {
    properties: {
      ".properties.diego_log_timestamp_format": {"value": "unix-epoch"},
    },
  };

  original_hash_rfc3339 = {
    properties: {
      ".properties.diego_log_timestamp_format": {"value": "rfc3339"},
    },
  };

  migrated_hash_rfc3339 = {
    properties: {
      ".properties.logging_timestamp_format": {"value": "rfc3339"},
    },
  };

  migrated_hash_deprecated = {
    properties: {
      ".properties.logging_timestamp_format": { "value": "deprecated"},
    },
  };

  describe("when the diego timestamp property is enabled", function(){
    it("enables the RFC3339 timestamp format", function() {
      migration.migrate(original_hash_rfc3339).should.deep.equal(migrated_hash_rfc3339);
    });
  });

  describe("when the diego timestamp property is not enabled", function(){
    it("does NOT enable the RFC3339 timestamp format", function() {
      migration.migrate(original_hash_epoch).should.deep.equal(migrated_hash_deprecated);
    });
  });
});

