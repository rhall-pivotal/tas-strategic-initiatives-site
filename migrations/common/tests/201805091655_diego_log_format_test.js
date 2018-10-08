require("tap").mochaGlobals()
const should = require("chai").should()
const migration = require("../201805091655_diego_log_format.js")

describe("Diego log timestamp format", function() {

  it("sets the value 'unix-epoch' on upgrade", function(){
    migration.migrate(
      { properties: {} }
    ).should.deep.equal(
      { properties: { ".properties.diego_log_timestamp_format": { "value": "unix-epoch" } } }
    );
  });
});
