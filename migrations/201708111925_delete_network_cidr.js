exports.migrate = function(input) {
    delete input.properties['.properties.container_networking.enable.network_cidr'];
    return input;
};
