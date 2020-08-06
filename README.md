# Application Service for Windows

![ballmer](http://i.giphy.com/mwDdHHKyuHe6c.gif)

## What it is?

Deploys garden cells (running windows server 2016) for all of your windows app needs

## Is there CI?

Yup, it is located [here](https://releng.ci.cf-app.com/?search=wrt)

## What you need to deploy it

- VMware Tanzu Application Service
- Windows Specific Stemcell

## How we test it

- Running [WATs](https://github.com/cloudfoundry/wats)
- Job currently configured with `skip_cert_verify`, if you don't do this the tests above will fail and the Rep won't be able to talk to Garden

## Contributing
asAaAS
Please refer to [PAS Contribution document](https://github.com/pivotal-cf/p-runtime/blob/master/CONTRIBUTING.md)

