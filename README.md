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
