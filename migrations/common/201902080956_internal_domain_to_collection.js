exports.migrate = function(input) {
  var properties = input.properties;

  properties[".properties.cf_networking_internal_domains"] = {
    value: [
      {
        guid: { value: generateGuid() },
        name: properties['.properties.cf_networking_internal_domain']
      }
    ]
  };

  delete properties['.properties.cf_networking_internal_domain'];

  return input;
};
