exports.migrate = function(input) {
  const properties = input.properties;

  if (properties[".properties.credhub_key_encryption_passwords"]) {
    let oldKeys =
      properties[".properties.credhub_key_encryption_passwords"].value;
    let internalKeys = [];
    let hsmKeys = [];
    oldKeys.map(function(key) {
      if (key.provider.value === "internal") {
        delete key.provider;
        internalKeys.push(key);
      } else {
        delete key.provider;
        hsmKeys.push(key);
      }
    });
    properties[".properties.credhub_internal_provider_keys"] = {
      value: internalKeys
    };
    properties[".properties.credhub_hsm_provider_encryption_keys"] = {
      value: hsmKeys
    };

    delete properties[".properties.credhub_key_encryption_passwords"];
  }
  return input;
};
