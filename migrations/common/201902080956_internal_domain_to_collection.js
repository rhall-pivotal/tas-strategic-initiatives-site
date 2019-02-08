exports.migrate = function(input) {
  var properties = input.properties;

  properties[".properties.cf_networking_internal_domains"] =
    [{ "name": properties['.properties.cf_networking_internal_domain']['value'] }];

  delete properties['.properties.cf_networking_internal_domain'];

  return input;
};
