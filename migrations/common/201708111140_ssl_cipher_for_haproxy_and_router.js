exports.migrate = function(input) {
  var properties = input.properties;

  var rfcToSSL = [
    ["", ""],
    ["TLS_RSA_WITH_RC4_128_SHA", "RC4-SHA"],
    ["TLS_RSA_WITH_3DES_EDE_CBC_SHA", "DES-CBC3-SHA"],
    ["TLS_RSA_WITH_AES_128_CBC_SHA", "AES128-SHA"],
    ["TLS_RSA_WITH_AES_256_CBC_SHA", "AES256-SHA"],
    ["TLS_RSA_WITH_AES_128_GCM_SHA256", "AES256-GCM-SHA384"],
    ["TLS_RSA_WITH_AES_256_GCM_SHA384", "AES256-GCM-SHA384"],
    ["TLS_ECDHE_ECDSA_WITH_RC4_128_SHA", "ECDHE-ECDSA-RC4-SHA"],
    ["TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA", "ECDHE-ECDSA-AES128-SHA"],
    ["TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA", "ECDHE-ECDSA-AES256-SHA"],
    ["TLS_ECDHE_RSA_WITH_RC4_128_SHA", "ECDHE-RSA-RC4-SHA"],
    ["TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA", "ECDHE-RSA-DES-CBC3-SHA"],
    ["TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA", "ECDHE-RSA-AES128-SHA"],
    ["TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA", "ECDHE-RSA-AES256-SHA"],
    ["TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "ECDHE-RSA-AES128-GCM-SHA256"],
    ["TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256", "ECDHE-ECDSA-AES128-GCM-SHA256"],
    ["TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", "ECDHE-RSA-AES256-GCM-SHA384"],
    ["TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "ECDHE-ECDSA-AES256-GCM-SHA384"],
  ];

  var cipherMap = new Map(rfcToSSL);

  if( properties['.properties.networking_point_of_entry'] ) {
    if( properties['.properties.networking_point_of_entry']['value'] == 'external_ssl' ) {
      if( properties['.properties.networking_point_of_entry.external_ssl.ssl_ciphers']['value'] != null ) {
        var rfcCiphers = properties['.properties.networking_point_of_entry.external_ssl.ssl_ciphers']['value'].split(':')
        var sslCiphers = rfcCiphers.map(function(c) {
          return cipherMap.get(c)
        })

        properties['.properties.gorouter_ssl_ciphers'] = {
          value: sslCiphers.join(':')
        }
      }
    } else if ( properties['.properties.networking_point_of_entry']['value'] == 'haproxy' ) {
      properties['.properties.haproxy_ssl_ciphers'] = properties['.properties.networking_point_of_entry.haproxy.ssl_ciphers'];
    }
  }

  return input;
};
