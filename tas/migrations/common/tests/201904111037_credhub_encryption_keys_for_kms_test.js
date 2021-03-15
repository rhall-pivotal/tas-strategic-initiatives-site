require("./spec_helper.js");
const migration = require("../201904111037_credhub_encryption_keys_for_kms.js");

describe("credhub encryption key migrations", function() {
  context("when there are only internal keys", function() {
    it("migrates internal keys to internal_provider_keys", function() {
      migration
        .migrate({
          properties: {
            ".properties.credhub_key_encryption_passwords": {
              value: [
                {
                  name: { value: "Key" },
                  key: { value: "some-long-internal-key" },
                  provider: { value: "internal" },
                  primary: { value: true }
                }
              ]
            }
          }
        })
        .should.deep.equal({
          properties: {
            ".properties.credhub_internal_provider_keys": {
              value: [
                {
                  name: { value: "Key" },
                  key: { value: "some-long-internal-key" },
                  primary: { value: true }
                }
              ]
            },
            ".properties.credhub_hsm_provider_encryption_keys": {
              value: []
            }
          }
        });
    });
  });

  context("when there are internal and hsm keys", function() {
    it("migrates keys to respective collections", function() {
      migration
        .migrate({
          properties: {
            ".properties.credhub_key_encryption_passwords": {
              value: [
                {
                  name: { value: "Key" },
                  key: { value: "some-long-internal-key" },
                  provider: { value: "internal" },
                  primary: { value: false }
                },
                {
                  name: { value: "Key" },
                  key: { value: "some-other-long-internal-key" },
                  provider: { value: "internal" },
                  primary: { value: false }
                },
                {
                  name: { value: "Key" },
                  key: { value: "some-long-hsm-provider-key" },
                  provider: { value: "hsm" },
                  primary: { value: false }
                },
                {
                  name: { value: "Key" },
                  key: { value: "some-other-long-hsm-provider-key" },
                  provider: { value: "hsm" },
                  primary: { value: true }
                }
              ]
            }
          }
        })
        .should.deep.equal({
          properties: {
            ".properties.credhub_internal_provider_keys": {
              value: [
                {
                  name: { value: "Key" },
                  key: { value: "some-long-internal-key" },
                  primary: { value: false }
                },
                {
                  name: { value: "Key" },
                  key: { value: "some-other-long-internal-key" },
                  primary: { value: false }
                }
              ]
            },
            ".properties.credhub_hsm_provider_encryption_keys": {
              value: [
                {
                  name: { value: "Key" },
                  key: { value: "some-long-hsm-provider-key" },
                  primary: { value: false }
                },
                {
                  name: { value: "Key" },
                  key: { value: "some-other-long-hsm-provider-key" },
                  primary: { value: true }
                }
              ]
            }
          }
        });
    });
  });
});
