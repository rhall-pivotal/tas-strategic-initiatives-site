require("./spec_helper.js");

const migration = require("../202002140809_routing_sticky_session_cookies.js");

describe("GoRouter sticky session cookies", function() {
  it("defaults sticky session cookies to JSESSIONID", function() {
    var migratedProperties = migration.migrate({ properties: {} })
    var sessionCookies = migratedProperties['properties']['.properties.router_sticky_session_cookie_names']['value']

    sessionCookies.length.should.equal(1)
    sessionCookies[0]['name']['value'].should.equal("JSESSIONID");
  });
});

