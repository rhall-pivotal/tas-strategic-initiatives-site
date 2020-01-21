exports.migrate = function(input) {
  input.properties['.properties.router_sticky_session_cookie_names'] = {
    "value": [
      {
        guid: { value: generateGuid() },
        name: { value: 'JSESSIONID' }
      }
    ]
  };

  return input;
};

