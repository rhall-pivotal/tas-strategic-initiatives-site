# p-runtime

This repository is used to generate a .pivotal file for Pivotal's Elastic
Runtime (see [Creating a Pivotal Cloud Foundry Product
Tile](https://docs.pivotal.io/partners/creating.html)), to be consumed by
Operations Manager&trade;.

## Updating an Instance Group

Runtime definition for an instance group can be found in
`instance_groups/<instance_group_name>.yml`, and configuration (forms and
validations) can be found in `forms/<job_name>.yml`. See the [Product Template
Reference](https://docs.pivotal.io/partners/product-template-reference.html)
for details on the formats.

### Changing Instance Group Order

Instance group order is defined by `_order.yml` files. To change the order
instance groups appear in the Ops Manager UI, edit `forms/_order.yml`. To
change the order in which the instance groups are deployed, edit
`instance_groups/_order.yml`

## Contributing

p-runtime is used to build .pivotal files for all supported versions of PCF,
as well as future versions. At the time of writing, this is PCF
[1.10.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.10),
[1.11.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.11),
[1.12.x](https://github.com/pivotal-cf/p-runtime/tree/rel/1.12),
and [2.0.x](https://github.com/pivotal-cf/p-runtime/tree/rel/2.0).
Each version is represented as a branch and must be updated independently.
If a change is required in more than one version, separate PRs for each
branch will be required.
