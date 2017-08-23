# p-windows-runtime-2016

![ballmer](http://i.giphy.com/mwDdHHKyuHe6c.gif)

## What it is?

Deploys garden cells (running windows server 2016) for all of your windows app needs

## Is there CI?

TBA

## What you need to deploy it

- Elastic Runtime
- Windows Specific Stemcell

## How we test it

- Running [WATs](https://github.com/cloudfoundry/wats)
- Job current configured with `skip_cert_verify`, if you don't do this the tests above will fail and the Rep won't be able to talk to Garden

## What's the catch???

- SSH is problematic (you can't for now). You need to Remote Desktop if you want to get on the Cells
  - To allow remote desktop connectivity you have to edit a registry setting
