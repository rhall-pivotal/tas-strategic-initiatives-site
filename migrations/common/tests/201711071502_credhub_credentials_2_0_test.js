require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201711071502_credhub_credentials_2_0.js");

original_hash = {
	properties: {
		".properties.deploy_autoscaling_broker_credentials": "super-secret-1",
		".properties.deploy_autoscaling_encryption_key": "super-secret-2",
		".backup-prepare.backup_encryption_key": "super-secret-3",
		".diego_database.bbs_encryption_passphrase": "super-secret-4",
		".nfs_server.blobstore_secret": "super-secret-5",
		".properties.deploy_notifications_encryption_key": "super-secret-6",
		".properties.push_pivotal_account_encryption_key": "super-secret-7",
		".properties.push_usage_service_secret_token": "super-secret-8",
		".router.route_services_secret": "super-secret-9",
	},
  variable_migrations: []
};

migrated_hash = {
	properties: {
		".properties.deploy_autoscaling_broker_credentials": "super-secret-1",
		".properties.deploy_autoscaling_encryption_key": "super-secret-2",
		".backup-prepare.backup_encryption_key": "super-secret-3",
		".diego_database.bbs_encryption_passphrase": "super-secret-4",
		".nfs_server.blobstore_secret": "super-secret-5",
		".properties.deploy_notifications_encryption_key": "super-secret-6",
		".properties.push_pivotal_account_encryption_key": "super-secret-7",
		".properties.push_usage_service_secret_token": "super-secret-8",
		".router.route_services_secret": "super-secret-9",
	},
  variable_migrations: [
    {
      "from": "super-secret-1",
      "to_variable": "deploy-autoscaling-broker-credentials"
    },
    {
      "from": "super-secret-2",
      "to_variable": "deploy-autoscaling-encryption-key"
    },
    {
      "from": "super-secret-3",
      "to_variable": "deploy-autoscaling-encryption-key"
    },
    {
      "from": "super-secret-4",
      "to_variable": "diego-db-bbs-encryption-passphrase"
    },
    {
      "from": "super-secret-5",
      "to_variable": "nfs-server-blobstore-secret"
    },
    {
      "from": "super-secret-6",
      "to_variable": "deploy-notifications-encryption-key"
    },
    {
      "from": "super-secret-7",
      "to_variable": "push-pivotal-account-encryption-key"
    },
    {
      "from": "super-secret-8",
      "to_variable": "push-usage-service-secret-token"
    },
    {
      "from": "super-secret-9",
      "to_variable": "router-route-services-secret"
    },
  ]
};

describe("migrate 2.0 credentials to credhub", function() {
  context("when migrations run", function() {
    it("returns the expected migrated hash", function(){
      migration.migrate(original_hash).should.deepEqual(migrated_hash);
    });
  });
});
