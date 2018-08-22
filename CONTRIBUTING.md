## Contributing changes to PAS

**Note:** this doc is a work-in-progress

p-runtime is used to build `.pivotal` files for all supported versions of PCF,
as well as future versions.
Each version is represented as a branch, e.g. [rel/2.2](https://github.com/pivotal-cf/p-runtime/tree/rel/2.0), and must be updated independently.
The `master` branch represents the "next" version of PAS and is used to build
release-candidates for the upcoming release.
If a change is required in more than one version, separate PRs for each branch will be required.

#### Building a tile

To build a PAS tile locally to test your changes:
1. Download an install the latest
   [kiln](https://github.com/pivotal-cf/kiln/releases) binary
1. Make a `./releases` directory in the p-runtime repo
1. Download all the BOSH releases into `./releases`
  - This is a manual step for the moment, we usually download a previously built
    tile and unzip the `./releases` dir from it
1. To build a Small Footprint PAS run `./bin/build` or `PRODUCT=ert ./bin/build` to
   build a full PAS tile
1. If you only need to test UI changes and don't need to actually deploy the
   tile, you can skip the release downloading with `STUB_RELEASES=true
   ./bin/build`

#### Property Changes

To TDD changes to a PAS tile, use [planitest](https://github.com/pivotal-cf/planitest) and [ops-manifest](https://github.com/pivotal-cf/ops-manifest). `planitest` is a testing library to make assertions against the generated BOSH manifest while `ops-manifest` is extracted Ops Manager code that transforms the tile metadata into a BOSH manifest without the need to stand up a running Ops Manager. 

When you write tests using `planitest` and `ops-manifest`, you are testing that changes made to the tile will result in the expected changes to the BOSH manifest, such as adding a job to an instance group or ensuring a property value is set. Currently, there is no way to test changes to the tile UI. 

[See example code here](https://github.com/pivotal-cf/planitest/blob/master/example_product_service_test.go).

If you want more confidence in your tile changes, you can use [om](https://github.com/pivotal-cf/om/) as the renderer for `planitest` and run your tests against a real Ops Manager.

For more information about tile metadata, refer to the [Product Template Reference](https://docs.pivotal.io/tiledev/2-2/product-template-reference.html).

**Setting up ops-manifest and running tests**

1. Clone the branch of PAS you want to make changes to
1. Generate a Github API token
1. Run `./bin/test` with your Github username and API token
1. To add new tests, modify or create an appropriate test file in `tests/manifest`
1. Implement the code changes in the PAS repo, IST, WRT, etc.
1. Make a PR!

#### Bumping releases

The current process for getting a BOSH release updated in PAS is to open an
issue on the [lts-pas-issues repo](https://github.com/pivotal-cf/lts-pas-issues/issues) and the PAS RelEng team will make the change for you.
To see which version of your release is currently included in PAS, run this:
```
curl https://releng.ci.cf-app.com/api/v1/teams/main/pipelines/build::2.3/resources
```

We have an upcoming track of work to move towards a more self-service model for teams where the p-runtime repo will contain an `assets.yml` or similar which lists the release versions.
See more about this epic [here](https://www.pivotaltracker.com/epic/show/4007210).

#### Migrations

**What is a migration?**

Ops Manager allows tile authors to write JavaScript [migration
files](https://docs.pivotal.io/tiledev/2-2/tile-upgrades.html#import) in order
to modify the value of selected properties when upgrading from one version of a
tile to a newer version.
For example, with migrations you can:
- Introduce a new default value, but have existing environments keep the
  previous default
  - e.g. existing environments should have the Diego Logging format set to
    `unix-epoch` on upgrade while new environments should be set to `rfc3339`
- Prevent an Operator from upgrading unless certain properties were selected in
  the previous version
  - e.g. you must enable grootfs prior to upgrading to PAS 2.3, else raise an
    error
- Delete an unused property from the Ops Manager property database

Ops Manager will sort the migration files alphabetically by filename and apply
them in ascending order. On a clean install, all migrations present in the tile
will automatically be marked as "ran". On upgrade, any migration files that are
present in the new tile but were not present in the previous tile version will
be applied.

**How do you write a migration?**

Migration files are JavaScript files which take in a hash of existing
properties, modify those properties if necessary, then return the updated hash
of properties. The filename of the migration should start with the current
datetime to ensure this migration runs after all existing migrations.

For example we can create a migration file on the rel/2.2 branch called `migrations/common/201805091655_diego_log_format.js` with the following contents:

```js
exports.migrate = function(input) {
  var properties = input.properties;

  properties[".properties.diego_log_timestamp_format"] = { "value": "unix-epoch" };

  return input;
};
```

When the Operator upgrades from PAS 2.1 to PAS 2.2, this migration will set the
Diego Logging format to `unix-epoch` rather than the new install default of
`rfc3339` to maintain existing behavior.

To test-drive this change, you can create a test file similar to the following
with the name `migrations/common/tests/201805091655_diego_log_format_test.js`:

```js
require("tap").mochaGlobals()
const should = require("should")
const migration = require("../201805091655_diego_log_format.js")

describe("Diego log timestamp format", function() {

    it("sets the value 'unix-epoch' on upgrade", function(){
        migration.migrate(
            { properties: {} }
        ).should.deepEqual(
            { properties: { ".properties.diego_log_timestamp_format":
              { "value": "unix-epoch" } } }
        );
    });
});
```

To run the migration tests:

```
cd ./migrations
npm install
npm test
```

**Are there any gotchas?**

- As a general rule, **do NOT modify** migration files which have already shipped in a version of PAS.
  For example, if you wrote a migration that shipped in PAS 2.0.0 but modified
  that migration in PAS 2.0.1, then Operators who upgraded from 2.0.0 to 2.0.1
  get migration version A while Operators who upgraded from 1.12.0 to 2.0.1 would get migration version B.

  For example:
  1. You wrote a migration that shipped in PAS 2.0.0
  2. You modified that migration in PAS 2.0.1
  3. Operators who upgrade from 2.0.0 to 2.0.1 get migration version A
  4. Operators who upgrade from 1.12.0 to 2.0.1 get migration version B

  This sort of diverge schema can create subtle bugs and make handling support
  tickets more difficult.
- If you remove a property and its associated form fields from the tile, Ops
  Manager will still have that property and its value stored in its property
  database. Therefore it's good practice to write a migration to delete that
  property from the database.
- Normally properties are a nested hash with a `value` key which is easy to
  forget.

  Correct:
  ```js
  properties[".properties.diego_log_timestamp_format"] = { "value": "unix-epoch" };
  ```

  Wrong:
  ```js
  properties[".properties.diego_log_timestamp_format"] = "unix-epoch"; // missing { "value": ... }
  ```
