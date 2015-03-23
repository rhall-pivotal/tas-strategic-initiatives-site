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

### Deploying to Ops Manager


If you want to start fresh, you can clear out your environment completely:

**WARNING: this destroys all VMs in the _faux_ environment, including but not limited to Ops Manager, µBOSH and Runtime**

```
bundle exec rake opsmgr:destroy[faux]
```

Deploy the Ops Manager .ova:

```
bundle exec rake opsmgr[faux,/path/to/your/pivotal-vsphere-1.M.N.0.ova]
```

Once the Ops Manager is deployed, you can configure and deploy the µBOSH:

```
bundle exec rake opsmgr:bosh:redeploy[faux]
```

Once the µBOSH is deployed, you can configure and deploy the Elastic Runtime product:

```
bundle exec rake runtime[faux,/path/to/your/cf-1.X.Y.0.pivotal]
```
