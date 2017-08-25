exports.migrate = function(input) {
  input.variable_migrations.push({
    from: input.properties['.mysql.autoscale_credentials'],
    to_variable: 'mysql-autoscale-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.ccdb_credentials'],
    to_variable: 'mysql-cc-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.diag_agent_credentials'],
    to_variable: 'mysql-diag-agent-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.diegodb_credentials'],
    to_variable: 'mysql-diego-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.locketdb_credentials'],
    to_variable: 'mysql-locket-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.monitordb_credentials'],
    to_variable: 'mysql-monitor-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.mysql_backup_server_credentials'],
    to_variable: 'mysql-backup-server-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.mysql_bootstrap_credentials'],
    to_variable: 'mysql-bootstrap-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.networkpolicyserverdb_credentials'],
    to_variable: 'mysql-network-policy-server-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.nfsvolume_credentials'],
    to_variable: 'mysql-nfs-volume-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.notifications_credentials'],
    to_variable: 'mysql-notifications-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.pivotal_account_credentials'],
    to_variable: 'mysql-pivotal-account-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.routingdb_credentials'],
    to_variable: 'mysql-routing-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.silkdb_credentials'],
    to_variable: 'mysql-silk-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.mysql.uaadb_credentials'],
    to_variable: 'mysql-uaa-db-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.nfsbrokerpush.nfs_broker_push_credentials'],
    to_variable: 'nfs-broker-push-credentials'
  });
  input.variable_migrations.push({
    from: input.properties['.notifications.encryption_credentials'],
    to_variable: 'notifications-encryption-credentials'
  });
  return input;
};
