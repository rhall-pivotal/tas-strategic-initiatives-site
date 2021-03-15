exports.migrate = function(input) {
    var properties = input.properties;

    if( properties['.uaa.sso_name']['value'] ) {
        properties['.properties.uaa'] = {value: 'saml'};

        properties['.properties.uaa.saml.sso_name'] = properties['.uaa.sso_name'];
        properties['.properties.uaa.saml.display_name'] = properties['.uaa.sso_name'];
        properties['.properties.uaa.saml.sso_xml'] = properties['.uaa.sso_xml'];
        properties['.properties.uaa.saml.sso_url'] = properties['.uaa.sso_url'];
    } else if( properties['.uaa.ldap_url']['value'] ) {
        properties['.properties.uaa'] = {value: 'ldap'};

        properties['.properties.uaa.ldap.url'] = properties['.uaa.ldap_url'];
        properties['.properties.uaa.ldap.credentials'] = properties['.uaa.credentials'];
        properties['.properties.uaa.ldap.search_base'] = properties['.uaa.ldap_search_base'];
        properties['.properties.uaa.ldap.search_filter'] = properties['.uaa.ldap_search_filter'];
        properties['.properties.uaa.ldap.group_search_base'] = properties['.properties.group_search.enable_admin_groups.ldap_group_search_base'];
        properties['.properties.uaa.ldap.group_search_filter'] = properties['.properties.group_search.enable_admin_groups.ldap_group_search_filter'];
        properties['.properties.uaa.ldap.server_ssl_cert_alias'] = properties['.uaa.ldap_server_ssl_cert_alias'];
        properties['.properties.uaa.ldap.server_ssl_cert'] = properties['.uaa.ldap_server_ssl_cert'];
        properties['.properties.uaa.ldap.mail_attribute_name'] = properties['.uaa.ldap_mail_attribute_name'];
    } else {
        properties['.properties.uaa'] = {value: 'internal'};
    }

    return input;
};
