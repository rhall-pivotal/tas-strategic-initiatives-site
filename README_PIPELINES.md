## Feature Pipelines
Feature pipelines are created for feature branches. They can be either
clean-installs or upgrades, and can be created for any IAAS.

The templates are in `ci/pipelines/feature*template.yml`

The tooling also uses templates in `ci/pipelines/release/template/*`

If we want a feature branch to have multiple types of pipelines, we create
a pipeline, rename it, then create the next one.

## Batch Pipelines
Batch pipelines are created with the same command used to create feature
pipelines. Batch pipelines are usually upgrade pipelines.

## Release Pipelines
Release pipelines (1.6 and later) are created by rake tasks, and are built from templates
in `ci/pipelines/release/template/*`

There is a full-suite pipeline, and a half-suite pipeline.

## Edge Pipeline
The edge pipelines consume new ops manager builds, use the latest ert.pivotal
promoted by a release pipeline, and use the tooling on the master branch of
p-runtime.

## Salmon Pipelines
Salmon pipelines are co-owned with upstream teams, and run against salmon
branches augmented with an override file for binaries.yml.

The salmon/cf pipeline is maintained by releng.

## Util Pipelines
The build docker image pipeline runs once per day, and pushes a docker image
to dockerhub.
