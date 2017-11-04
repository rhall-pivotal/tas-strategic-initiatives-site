exports.migrate = function(input) {
    var properties = input.properties;

    properties['.cloud_controller.internal_api_user_credentials'] = properties['.health_manager.internal_api_user_credentials'];

    return input;
};
