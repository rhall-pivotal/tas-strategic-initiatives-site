## p-runtime

Pivotal's Elastic Runtime tile, to be consumed by Ops Manager&trade;.

Relies on a [fork](https://github.com/pivotal-cf/pcf-release) of Cloud Foundry's [open source elastic runtime](https://github.com/cloudfoundry/cf-release).

### Creating .pivotal file

```
git clone git@github.com:pivotal-cf/p-runtime.git
cd p-runtime
bundle install # Installs vara
bundle exec vara build-pivotal ~/workspace/p-runtime/ # Creates cf-1.N.0.0.alpha.XYZ.sOmEsHa.pivotal
```
### Configuring p-runtime / opsmgr tasks

#### Environments

A named environment is specified in `config/environments/<name>.yml` (or `${ENV_DIRECTORY)/<name>.yml` if that environment variable is set.

See [the p-runtime environment files](https://github.com/pivotal-cf/p-runtime/blob/master/config/environments) for examples of defining your environments.

### Deploying Ops Manager and Elastic Runtime

#### Clean existing environment
If you want to start fresh, you can clear out your environment completely:

**WARNING: this destroys all VMs in the environment, including but not limited to Ops Manager, µBOSH and Runtime**

Note: If you used the old version of opsmgr tasks to create this environment, then you need to use that version to clean it. The new version does not correctly clean out deployments from the old version.

```
bundle exec rake opsmgr:destroy[environment]
```

#### Download an Ops Manager Image
Obtain the ACCESS_KEY_ID and SECRET_ACCESS_KEY from Last Pass. The credentials are in `Shared - Ops Manager Secure` as `AWS "+opsmanager" Account "RelengToolsReader" IAM User access key`. Contact the Ops Manager team for access.

Replace \<iaas> with `openstack`, `aws`, or `vsphere`

```
ACCESS_KEY_ID=<OpsMan Reader Access Key> SECRET_ACCESS_KEY=<OpsMan Reader Secret Key> bundle exec rake opsmgr:download[<iaas>,stable]
```

#### Deploy the Ops Manager
The above task saves the Ops Manager Image to a file `ops_man_image`

```
bundle exec rake opsmgr:install[environment,ops_man_iamge]
```

#### Configure and Deploy Microbosh
You need to specify the major and minor version of Ops Manager in these commands.

`<OM version>` is the Ops Manager version. Opsmgr supports Ops Manager `1.4` and `1.5`.

`<wait time>` is number of minutes to wait for install, recommended wait time is `45`

```
bundle exec rake opsmgr:add_first_user[environment,<OM version>]
bundle exec rake opsmgr:microbosh:configure[environment,<OM version>]
bundle exec rake opsmgr:trigger_install[environment,<OM version>,<wait time>]
```

#### Upload, Configure, and Deploy ERT
Once the µBOSH is deployed, you can configure and deploy the Elastic Runtime product

`<OM version>` is the Ops Manager version. Opsmgr supports Ops Manager `1.4` and `1.5`.

`<p-runtime .pivotal file>` is the .pivotal file created above

`cf` is the p-runtime product name

`<ert version>` is the ERT version. The p-runtime rake task supports ERT versions `1.4` and `1.5`

`<wait time>` is number of minutes to wait for install, recommended wait time is `240`

```
bundle exec rake opsmgr:product:upload_add[environment,<OM version>,<p-runtime .pivotal file>,cf]
bundle exec rake ert:configure[environment,<ert version>,<OM version>]
bundle exec rake opsmgr:trigger_install[environment,<OM version>,<wait time>]
```

#### Upgrade ERT to new version
`<OM version>` is the Ops Manager version. Opsmgr supports Ops Manager `1.4` and `1.5`.

`<new p-runtime .pivotal file>` is the newer version of the .pivotal file created above

`cf` is the p-runtime product name

`<ert version>` is the ERT version. The p-runtime rake task supports ERT versions `1.4` and `1.5`

`<wait time>` is number of minutes to wait for install, recommended wait time is `240`

```
bundle exec rake opsmgr:product:upload_upgrade[environment,<OM version>,<new p-runtime .pivotal file>,cf]
bundle exec rake ert:configure[environment,<ert version>,<OM version>]
bundle exec rake opsmgr:trigger_install[environment,<OM version>,<wait time>]
```

#### Export the installation from Ops Manager
`<OM version>` is the Ops Manager version. Opsmgr supports Ops Manager `1.4` and `1.5`.

`<file name>` local file name to save the exported installation file

```
bundle exec rake opsmgr:export_installation[environment,<OM version>,<file name>]
```

#### Import an installation to Ops Manager
`<OM version>` is the Ops Manager version. Opsmgr supports Ops Manager `1.4` and `1.5`.

`<file name>` local file name of the installation file to import

```
bundle exec rake opsmgr:import_installation[environment,<OM version>,<file name>]
```

#### Destroy only the Ops Manager VM

Commonly done when testing export/import scenarios

```
bundle exec rake opsmgr:destroy:opsmgr[environment]
```

### ERT Tasks as examples for your Pivotal product

The [Opsmgr gem](https://github.com/pivotal-cf/opsmgr) brings in the [Ops Manager UI Drivers gem](https://github.com/pivotal-cf-experimental/ops_manager_ui_drivers) which can be used to create configuration tasks for your product. To create configuration tasks for your product, you should add the following code to your product:

1. Add the Opsmgr gem as a dependency to your product
1. Add a [section to your environment files](https://github.com/pivotal-cf/p-runtime/blob/master/config/environments/skunk.yml#L61-L75) with the configuration details
1. Create a [rake task](https://github.com/pivotal-cf/p-runtime/blob/master/lib/tasks/ert.rake) in your product
1. Create an [integration spec runner](https://github.com/pivotal-cf/p-runtime/blob/master/lib/tools/integration_spec_runner.rb) to invoke the correct integration test for your product version
1. Create [integration tests](https://github.com/pivotal-cf/p-runtime/tree/master/integration) that execute the desired configuration tasks
