exports.migrate = function(input) {
    delete input.properties['.properties.container_networking.disable.garden_network_pool'];
    return input;
};
