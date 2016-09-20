# p-runtime

This repository is used to generate a .pivotal file for Pivotal's Elastic Runtime (see [Creating a Pivotal Cloud Foundry Product Tile](https://docs.pivotal.io/partners/creating.html)), to be consumed by Ops Manager&trade;.

## Updating a Job

Runtime definition for a job can be found in `metadata_parts/jobs/<job_name>.yml`, and configuration (forms and validations) can be found in `metadata_parts/forms_and_validators/<job_name>.yml`. See the [Product Template Reference](https://docs.pivotal.io/partners/product-template-reference.html) for details on the formats.

### Changing Job Order

Job order is defined by `_order.yml` files. To change the order jobs appear in the Ops Manager UI, edit `metadata_parts/forms_and_validators/_order.yml`. To change the order in which the jobs are deployed, edit `metadata_parts/jobs/_order.yml`

## Contributing

p-runtime is used to build .pivotal files for all supported versions of PCF, as well as future versions. At the time of writing, this is PCF [1.6.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.6), [1.7.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.7), [1.8.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.8), and [1.9.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.9). Each version is represented as a branch and must be updated independently. If a change is required in more than one version, separate PRs for each branch will be required.
