## p-runtime

**p-runtime** creates .pivotal files to be consumed by Ops Manager&trade;.  It creates .pivotal files based on [Elastic Runtime](https://github.com/cloudfoundry/cf-release) releases.

### Creating the CF Release .tgz
[scripts/build_cf_release_tarball.sh](https://github.com/pivotal-cf/release-engineering-automation/blob/master/scripts/build_cf_release_tarball.sh) script in the [pivotal-cf/release-engineering-automation](https://github.com/pivotal-cf/release-engineering-automation) repo

```
cd pivotal-cf/release-engineering-automation
scripts/build_cf_release_tarball.sh
```

### Creating .pivotal File

In this example, we update the Elastic Runtime.  We assume the following:

* The product version is **1.3.0.0**
* The CF Release tag is **v168**
* The stemcell version is **2399**

We follow this procedure:

* Edit *metadata_parts/binaries.yml* and modify the entries to match the following:

```
stemcell:
  name: bosh-vsphere-esxi-ubuntu
  version: '2399'
  file: bosh-stemcell-2399-vsphere-esxi-ubuntu.tgz
  md5: ee144a0a3abed4cdf543cef11e9f0f7b
releases:
- file: cf-168.tgz
  name: cf
  version: '168'
  md5: 88ba12abac6e44cba0562e6c812daf4a
  url: https://releng-artifacts.s3.amazonaws.com/cf-168.tgz
```

* Edit *metadata/handcraft.yml* and update the product number to *1.3.0.0*.  Also, update the metadata version (used by Ops Manager) to *1.3*:

```
name: cf
product_version: 1.3.0.0
metadata_version: '1.3'
provides_product_versions:
  - name: cf
    version: 1.3.0.0
```

* Use [vara](https://github.com/pivotal-cf/vara) to re-generate the metadata .yml file by merging the metadata partial files:

```
vara-build-metadata --product-dir=~/workspace/p-runtime
```

* Use *vara* to download the stemcells and artifacts:

```
vara-download-artifacts --product-metadata=~/workspace/p-runtime/metadata/cf.yml
```

* Use *vara* (again) to create the .pivotal file:

```
vara-build-pivotal --product-metadata=~/workspace/p-runtime/metadata/cf.yml
```

* Look for file named *cf-168.pivotal*
* Upload that file to an Ops Manager VM
* Configure the settings
* Install
* Test (e.g. running [CATS](https://github.com/cloudfoundry/cf-acceptance-tests))

Any failures will mostly likely be addressed by modifying *metadata/handcraft.yml* and rebuilding the .pivotal file.