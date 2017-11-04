exports.migrate = function(input) {
    var properties = input.properties;

    if (properties['.properties.mysql_backups']['value'] == 'enable') {
        properties['.properties.mysql_backups']['value'] = 's3';

        properties['.properties.mysql_backups.s3.endpoint_url'] = properties['.properties.mysql_backups.enable.endpoint_url'];
        properties['.properties.mysql_backups.s3.bucket_name'] = properties['.properties.mysql_backups.enable.bucket_name'];
        properties['.properties.mysql_backups.s3.bucket_path'] = properties['.properties.mysql_backups.enable.bucket_path'];
        properties['.properties.mysql_backups.s3.access_key_id'] = properties['.properties.mysql_backups.enable.access_key_id'];
        properties['.properties.mysql_backups.s3.secret_access_key'] = properties['.properties.mysql_backups.enable.secret_access_key'];
        properties['.properties.mysql_backups.s3.cron_schedule'] = properties['.properties.mysql_backups.enable.cron_schedule'];
        properties['.properties.mysql_backups.s3.backup_all_masters'] = properties['.properties.mysql_backups.enable.backup_all_masters'];
    }

    return input;
};
