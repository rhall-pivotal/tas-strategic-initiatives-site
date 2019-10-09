require("./spec_helper.js");

const migration = require("../201908161156_rename_add_compute_isolation_and_remove_placement_tag.js");

describe("a cell block upgrading from an isolation segment tag", function() {
  it("migrates old property in a non-replicated isolation segment", function(){
    var input = {
      "properties": {
        ".isolated_diego_cell.placement_tag": {
          "value": "isosegtag"
        }
      }
    };

    var expected = {
      "properties": {
        ".properties.compute_isolation": {
          "value": "enabled"
        },
        ".properties.compute_isolation.enabled.isolation_segment_name": {
          "value": "isosegtag"
        }
      }
    };

    migration.migrate(input).should.deep.equal(expected);
  })

  it("migrates old property in a replicated isolation segment", function(){
    var input = {
      "properties": {
        ".isolated_diego_cell_replicated_foo.placement_tag": {
          "value": "isosegtag"
        }
      }
    };

    var expected = {
      "properties": {
        ".properties.compute_isolation": {
          "value": "enabled"
        },
        ".properties.compute_isolation.enabled.isolation_segment_name": {
          "value": "isosegtag"
        }
      }
    };

    migration.migrate(input).should.deep.equal(expected);
  })
});
