require("./spec_helper.js");

const migration = require("../201803091153_delete_mysql_backups_selector.js");

describe("Remove Automated Backup Configuration", function() {
  original_hash = {
    properties: {
      ".properties.mysql_backups": {
        "type": "selector",
        "configurable": true,
        "credential": false,
        "value": "disable",
        "optional": false
      },
      ".properties.mysql_backups.s3.endpoint_url": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": "https://s3.amazonaws.com",
        "optional": false
      },
      ".properties.mysql_backups.s3.bucket_name": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.s3.bucket_path": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.s3.access_key_id": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.s3.secret_access_key": {
        "type": "secret",
        "configurable": true,
        "credential": true,
        "value": {
          "secret": "***"
        },
        "optional": false
      },
      ".properties.mysql_backups.s3.cron_schedule": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.s3.backup_all_masters": {
        "type": "boolean",
        "configurable": true,
        "credential": false,
        "value": true,
        "optional": false
      },
      ".properties.mysql_backups.s3.region": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": true
      },
      ".properties.mysql_backups.gcs.service_account_json": {
        "type": "secret",
        "configurable": true,
        "credential": true,
        "value": {
          "secret": "***"
        },
        "optional": false
      },
      ".properties.mysql_backups.gcs.project_id": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.gcs.bucket_name": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.gcs.cron_schedule": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.gcs.backup_all_masters": {
        "type": "boolean",
        "configurable": true,
        "credential": false,
        "value": true,
        "optional": false
      },
      ".properties.mysql_backups.azure.storage_account": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.azure.storage_access_key": {
        "type": "secret",
        "configurable": true,
        "credential": true,
        "value": {
          "secret": "***"
        },
        "optional": false
      },
      ".properties.mysql_backups.azure.container": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.azure.path": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.azure.cron_schedule": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.azure.backup_all_masters": {
        "type": "boolean",
        "configurable": true,
        "credential": false,
        "value": true,
        "optional": false
      },
      ".properties.mysql_backups.scp.server": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.scp.port": {
        "type": "port",
        "configurable": true,
        "credential": false,
        "value": 22,
        "optional": false
      },
      ".properties.mysql_backups.scp.user": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.scp.key": {
        "type": "text",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.scp.destination": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.scp.cron_schedule": {
        "type": "string",
        "configurable": true,
        "credential": false,
        "value": null,
        "optional": false
      },
      ".properties.mysql_backups.scp.backup_all_masters": {
        "type": "boolean",
        "configurable": true,
        "credential": false,
        "value": true,
        "optional": false
      },
    },
  };

  migrated_hash = {
    properties: {},
  };

  it("removes mysql_backup", function(){
    migration.migrate(original_hash).should.deep.equal(migrated_hash);
  });

});
