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

  input.variable_migrations.push({
    from: input.properties['.properties.consul_encrypt_key'],
    to_variable: 'consul-encryption-key'
  });

  input.variable_migrations.push({
    from: input.properties['.diego_database.bbs_encryption_passphrase'],
    to_variable: 'diego-db-bbs-encryption-passphrase'
  });

  input.variable_migrations.push({
    from: input.properties['.nats.credentials'],
    to_variable: 'nats-credentials'
  });

  input.variable_migrations.push({
    from: input.properties['.nfs_server.blobstore_secret'],
    to_variable: 'nfs-server-blobstore-secret'
  });

  input.variable_migrations.push({
    from: input.properties['.properties.deploy_notifications_encryption_key'],
    to_variable: 'deploy-notifications-encryption-key'
  });

  input.variable_migrations.push({
    from: input.properties['.properties.push_pivotal_account_encryption_key'],
    to_variable: 'push-pivotal-account-encryption-key'
  });

  input.variable_migrations.push({
    from: input.properties['.properties.push_usage_service_secret_token'],
    to_variable: 'push-usage-service-secret-token'
  });

  input.variable_migrations.push({
    from: input.properties['.router.route_services_secret'],
    to_variable: 'router-route-services-secret'
  });

  return input;
};
