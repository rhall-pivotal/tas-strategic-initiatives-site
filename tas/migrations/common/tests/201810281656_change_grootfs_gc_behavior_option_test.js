require("./spec_helper.js");

const migration = require("../201810281656_change_grootfs_gc_behavior_option.js");

describe("Grootfs GC deprecated graph_cleanup_threshold_in_mb for reserved_space_for_other_jobs_in_mb", function() {
  it("migrates threshold to reserved if default not changed", function(){
    var input = {
      "properties": {
        ".properties.garden_disk_cleanup": {
          "value": "threshold"
        },
        ".properties.garden_disk_cleanup.threshold.cleanup_threshold_in_mb": {
          "value": 10240
        }
      }
    };

    var expected = {
      "properties": {
        ".properties.garden_disk_cleanup": {
          "value": "reserved"
        },
        ".properties.garden_disk_cleanup.reserved.reserved_space_for_other_jobs_in_mb": {
          "value": 15360
        }
      }
    };

    migration.migrate(input).should.deep.equal(expected);
  })

  it("does not set a reserved property if threshold set to a non-default value", function(){
    var input = {
      "properties": {
        ".properties.garden_disk_cleanup": {
          "value": "threshold"
        },
        ".properties.garden_disk_cleanup.threshold.cleanup_threshold_in_mb": {
          "value": 20241
        }
      }
    };

    var expected = {
      "properties": {
        ".properties.garden_disk_cleanup": {
          "value": "threshold"
        },
      }
    };

    migration.migrate(input).should.deep.equal(expected);
  })
});
