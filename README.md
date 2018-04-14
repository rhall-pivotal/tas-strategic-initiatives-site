# p-runtime

This repository is used to generate a .pivotal file for Pivotal's Application
Service (see [Creating a Pivotal Cloud Foundry Product
Tile](https://docs.pivotal.io/partners/creating.html)), to be consumed by
Operations Manager&trade;.

## Building a PAS tile with this repo
The CLI tool [kiln](https://github.com/pivotal-cf/kiln) is used to build the .pivotal file.
See [kiln's Docs](https://github.com/pivotal-cf/kiln/blob/master/README.md) for more details.

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
