exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.uaa.pivotal_account_client_credentials'];

  return input;
};
