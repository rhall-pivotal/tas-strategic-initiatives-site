exports.migrate = function(input) {
    var properties = input.properties;

    properties['.properties.smoke_tests'] = {
        value: 'on_demand'
    };

    return input;
};
