require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201805091420_autoscale_metric_bucket.js")

describe("Autoscale Metric Bucket property increase", function() {

  it("sets the value to the default of 120", function(){
    migration.migrate(
      { properties: { ".properties.autoscale_metric_bucket_count": { "value": 35 } } }
    ).should.deepEqual(
      { properties: { ".properties.autoscale_metric_bucket_count": { "value": 120 } } }
    );
  });
});
