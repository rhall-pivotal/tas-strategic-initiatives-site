exports.migrate = function(input) {
  var properties = input.properties;

  properties['.properties.push_usage_service_secret_token'] = properties['.push-usage-service.secret_token'];

  properties['.properties.push_apps_manager_company_name'] = properties['.push-apps-manager.company_name'];
  properties['.properties.push_apps_manager_accent_color'] = properties['.push-apps-manager.accent_color'];
  properties['.properties.push_apps_manager_global_wrapper_bg_color'] = properties['.push-apps-manager.global_wrapper_bg_color'];
  properties['.properties.push_apps_manager_global_wrapper_text_color'] = properties['.push-apps-manager.global_wrapper_text_color'];
  properties['.properties.push_apps_manager_global_wrapper_header_content'] = properties['.push-apps-manager.global_wrapper_header_content'];
  properties['.properties.push_apps_manager_global_wrapper_footer_content'] = properties['.push-apps-manager.global_wrapper_footer_content'];
  properties['.properties.push_apps_manager_logo'] = properties['.push-apps-manager.logo'];
  properties['.properties.push_apps_manager_square_logo'] = properties['.push-apps-manager.square_logo'];
  properties['.properties.push_apps_manager_footer_text'] = properties['.push-apps-manager.footer_text'];
  properties['.properties.push_apps_manager_footer_links'] = properties['.push-apps-manager.footer_links'];
  properties['.properties.push_apps_manager_nav_links'] = properties['.push-apps-manager.nav_links'];
  properties['.properties.push_apps_manager_product_name'] = properties['.push-apps-manager.product_name'];
  properties['.properties.push_apps_manager_marketplace_name'] = properties['.push-apps-manager.marketplace_name'];
  properties['.properties.push_apps_manager_enable_invitations'] = properties['.push-apps-manager.enable_invitations'];
  properties['.properties.push_apps_manager_display_plan_prices'] = properties['.push-apps-manager.display_plan_prices'];
  properties['.properties.push_apps_manager_currency_lookup'] = properties['.push-apps-manager.currency_lookup'];

  return input;
};
