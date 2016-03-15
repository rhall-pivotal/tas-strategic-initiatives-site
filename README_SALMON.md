# Salmon Contract

"Salmon" teams ("teams", "you", "our upstream friends") that need to integrate their release into the Elastic Runtime tile ("ERT", "tile") can use this repository to build tiles including their releases. Once tested and delivered in a particular way, these releases can be consumed by the Release Engineering team ("RelEng", "we") for inclusion in the official tile. RelEng doesn't provide any particular pipelines to accomplish such building and delivering, but we have provided an example that teams can use as a guide [here]().

## Expectations Around Delivering Releases
We need:

- a publicly accessible tarball of the bosh release (either on S3 or github, for instance)
- the information necessary to update binaries.yml to use the new release, which looks like this:

```
---
- name: push-apps-manager-release
  file: push-apps-manager-release-452.tgz
  version: '452'
  md5: 874a3ce7b712c098cca8d8bfc75e6433
  url: http://apps-manager-releases.s3.amazonaws.com/push-apps-manager-release-452.tgz
```

## Expectations Around Delivering Changes to p-runtime
Each upstream release has a salmon branch on the p-runtime repo, which should be used to facilitate the building of modified tiles. Typically this will mean metadata changes. Regardless, requests to pull in changes from these branches should be made PM-to-PM, or through tracker, and should specify a branch to be merged, or particular shas to be cherry-picked.

## Expectations Around Testing Integrations
When a team delivers a release to RelEng, we expect that you will have:

- built an ERT including your release and any p-runtime changes you desire
- deployed the resulting tile to an Ops Manager of the current generation running on an IaaS of your choice
- run the Cloud Foundry Acceptance Tests ("CATs") against the resulting PCF deployment
- tested your release within that deployment to your satisfaction

We would also like it if teams have executable acceptance specs as part of the release (and run these as part of testing the integration), but this is not currently required. Typically this is done as an errand; CATs plays this role for cf-release, for instance.

## Help, Examples and Encouragement
RelEng is generally available to support the efforts of our upstream friends, but it is important that salmon pipelines be owned by you.