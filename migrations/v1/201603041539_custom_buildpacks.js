exports.migrate = function(input) {
  var properties = input.properties;

  // Switch from opt-out to opt-in language around custom buildpacks; rename the property and invert the value.
  properties['.cloud_controller.enable_custom_buildpacks'] = {
    value: ! properties['.cloud_controller.disable_custom_buildpacks']['value']
  };

  return input;
};
