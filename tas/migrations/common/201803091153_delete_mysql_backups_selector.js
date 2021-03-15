exports.migrate = function(input) {
  var properties = input.properties;

  delete properties['.properties.mysql_backups.s3.endpoint_url'];
  delete properties['.properties.mysql_backups.s3.bucket_name'];
  delete properties['.properties.mysql_backups.s3.bucket_path'];
  delete properties['.properties.mysql_backups.s3.region'];
  delete properties['.properties.mysql_backups.s3.access_key_id'];
  delete properties['.properties.mysql_backups.s3.secret_access_key'];
  delete properties['.properties.mysql_backups.s3.cron_schedule'];
  delete properties['.properties.mysql_backups.s3.backup_all_masters'];
  delete properties['.properties.mysql_backups.scp.server'];
  delete properties['.properties.mysql_backups.scp.port'];
  delete properties['.properties.mysql_backups.scp.user'];
  delete properties['.properties.mysql_backups.scp.key'];
  delete properties['.properties.mysql_backups.scp.destination'];
  delete properties['.properties.mysql_backups.scp.cron_schedule'];
  delete properties['.properties.mysql_backups.scp.backup_all_masters'];
  delete properties['.properties.mysql_backups.gcs.service_account_json'];
  delete properties['.properties.mysql_backups.gcs.project_id'];
  delete properties['.properties.mysql_backups.gcs.bucket_name'];
  delete properties['.properties.mysql_backups.gcs.cron_schedule'];
  delete properties['.properties.mysql_backups.gcs.backup_all_masters'];
  delete properties['.properties.mysql_backups.azure.storage_account'];
  delete properties['.properties.mysql_backups.azure.storage_access_key'];
  delete properties['.properties.mysql_backups.azure.container'];
  delete properties['.properties.mysql_backups.azure.path'];
  delete properties['.properties.mysql_backups.azure.cron_schedule'];
  delete properties['.properties.mysql_backups.azure.backup_all_masters'];
  delete properties['.properties.mysql_backups'];

  return input;
};
