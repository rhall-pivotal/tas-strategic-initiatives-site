# planitest

## What is it?

Test helpers for Ops Manager tile developers. Given the set of tile options selected by the operator, what should the generated BOSH manifest look like?

It can be prohibitively expensive to deploy your tile in each of these configurations - planitest lets you make assertions about the staged manifest.

## Example code

Basic assertion about the properties under an instance group. The `Property` method accepts a path expression (as you might use in ops-files):

```go
It("configures the router minimum TLS version", func() {
  err := product.Configure(map[string]interface{}{})
  Expect(err).NotTo(HaveOccurred())

  manifest, err := product.RenderManifest()
  Expect(err).NotTo(HaveOccurred())

  router, err := manifest.FindInstanceGroupJob("router", "gorouter")
  Expect(err).NotTo(HaveOccurred())
  Expect(router.Property("router/min_tls_version")).To(Equal("TLSv1.2"))
})
```

An example of a context that sets a different configuration. Here we override the default `routing_minimum_tls_version` and assert that the router is configured correctly:

```go
Context("when the operator sets the minimum TLS version to 1.1", func() {

  var manifest planitest.Manifest

  BeforeEach(func() {
    err := product.Configure(map[string]interface{}{
      ".properties.routing_minimum_tls_version": "tls_v1_1",
    })
    Expect(err).NotTo(HaveOccurred())

    manifest, err = product.RenderManifest()
    Expect(err).NotTo(HaveOccurred())
  })

  It("configures the router minimum TLS version", func() {
    router, err := manifest.FindInstanceGroupJob("router", "gorouter")
    Expect(err).NotTo(HaveOccurred())
    Expect(router.Property("router/min_tls_version")).To(Equal("TLSv1.1"))
  })
})
```

## What do you need?

1. An [Ops Manager](https://docs.pivotal.io/pivotalcf/1-12/customizing/) instance to test against. It should have the BOSH tile deployed.
1. The [om](https://github.com/pivotal-cf/om) CLI
1. The [bosh](https://bosh.io/docs/cli-v2.html#install) CLI
1. A minimal product-properties JSON file usable by `om configure-product`
1. A product-network JSON file usable by `om configure-product`
1. The tile you want to test. It should be already uploaded to Ops Manager, along with the stemcell it depends on.

## Rough edges

1. Don't attempt to run tests that use planitest in parallel as different examples will step on each other
1. Rendering a staged manifest for a large product on Ops Manager can be slooooow
1. Currently runs om with the `--skip-ssl-validation` flag
1. API is liable to change in breaking ways

## Prior art

* [om-manifest-validator](https://github.com/pivotal-cf-experimental/om-manifest-validator)
