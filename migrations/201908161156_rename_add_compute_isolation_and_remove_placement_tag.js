exports.migrate = function(input) {
    var properties = input['properties'];

    const keys = Object.keys(properties)

    const key = keys.find(k => /^\.isolated_diego_cell.*?placement_tag/.test(k))

    properties['.properties.compute_isolation.enabled.isolation_segment_name'] = { 'value': properties[key].value };
    properties['.properties.compute_isolation'] = {'value': 'enabled'};

    delete properties[key];

    return input;
};
