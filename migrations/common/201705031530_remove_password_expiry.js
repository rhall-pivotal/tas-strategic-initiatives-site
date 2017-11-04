exports.migrate = function(input) {
    var properties = input.properties;

    properties['.properties.uaa.internal.password_expires_after_months']['value'] = 0;

    return input;
};
