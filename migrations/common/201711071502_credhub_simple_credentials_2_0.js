exports.migrate = function(input) {

  input.variable_migrations.push({
    from: input.properties['.properties.deploy_autoscaling_broker_credentials'],
    to_variable: 'deploy-autoscaling-broker-credentials'
  });

  input.variable_migrations.push({
    from: input.properties['.properties.deploy_autoscaling_encryption_key'],
    to_variable: 'deploy-autoscaling-encryption-key'
  });

  input.variable_migrations.push({
    from: input.properties['.backup-prepare.backup_encryption_key'],
    to_variable: 'deploy-autoscaling-encryption-key'
  });

  return input;
};
