exports.migrate = function(input) {
    var properties = input['properties'];

    properties['.properties.compute_isolation'] = {'value': 'enabled'};
    properties['.properties.compute_isolation.enabled.isolation_segment_name'] = { 'value': properties['.isolated_diego_cell.placement_tag'].value };

    delete properties['.isolated_diego_cell.placement_tag'];

    return input;
};
