## p-runtime

**p-runtime** creates .pivotal files to be consumed by Ops Manager&trade;.  It creates .pivotal files based on [Elastic Runtime](https://github.com/cloudfoundry/cf-release) releases.

### Creating the CF Release .tgz
[scripts/build_cf_release_tarball.sh](https://github.com/pivotal-cf/release-engineering-automation/blob/master/scripts/build_cf_release_tarball.sh) script in the [pivotal-cf/release-engineering-automation](https://github.com/pivotal-cf/release-engineering-automation) repo

```
cd pivotal-cf/release-engineering-automation
scripts/build_cf_release_tarball.sh
```

### Creating .pivotal File (no compiled packages)

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

### Creating .pivotal File (with compiled packages)

We make sure we have successfully deployed a compiled packages-less .pivotal (we can't download compiled packages without first having deployed).

1. deploy a regular (no compiled packages) Elastic Runtime
2. ssh into the &micro;BOSH and download the compiled packages (only necessary for &micro;BOSHes who haven't fixed the nginx bug that truncates downloads to 1GB).  The exact procedure to follow has *not* been documented.
3. the resulting file should be an amalgam of the Release and the Stemcell, e.g. `cf-170-bosh-vsphere-esxi-ubuntu-2366.tgz`
4. put the aforementioned file in `~/workspace/p-runtime/lvl_2_compiled/compiled_packages/`


```
vim ~/workspace/p-runtime/lvl_2_compiled/metadata_parts/binaries.yml
```

Add the following lines:


```
compiled_package:
  name: cf
  file: cf-170-bosh-vsphere-esxi-ubuntu-2366.tgz
  version: "170"
  md5: a3eefb2dd839254e111d8e87232d036c
  url: https://releng-artifacts.s3.amazonaws.com/cf-170-bosh-vsphere-esxi-ubuntu-2366.tgz
```

We build the product file:

```
be vara-build-metadata --product-dir=~/workspace/p-runtime/lvl_2_compiled
```

We build the .pivotal file:


```
lvl_2_compiled/scripts/build_pivotal.sh
```

### Process to create subdirectory for compiled packages

```
cd ~/workspace/p-runtime
mkdir -p lvl_2_compiled/{metadata,metadata_parts,scripts}
cd lvl_2_compiled
for file in releases stemcells compiled_packages content_migrations metadata_parts/handcraft.yml; do
  ln -s ../$file $file
done

```
