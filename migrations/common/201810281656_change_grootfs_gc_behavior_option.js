exports.migrate = function(input) {
  var properties = input['properties'];


  if(properties['.properties.garden_disk_cleanup']['value'] !== 'threshold') {
    delete properties['.properties.garden_disk_cleanup.never.graph_cleanup_threshold_in_mb'];
    delete properties['.properties.garden_disk_cleanup.routine.graph_cleanup_threshold_in_mb'];
    return input;
  }

  var threshold = properties['.properties.garden_disk_cleanup.threshold.cleanup_threshold_in_mb']['value'];
  if(threshold == 10240) {
    properties['.properties.garden_disk_cleanup'] = {'value': 'reserved'};
    properties['.properties.garden_disk_cleanup.reserved.reserved_space_for_other_jobs_in_mb'] = {'value': 15360};
  }

  delete properties['.properties.garden_disk_cleanup.threshold.cleanup_threshold_in_mb'];
  return input;
};
