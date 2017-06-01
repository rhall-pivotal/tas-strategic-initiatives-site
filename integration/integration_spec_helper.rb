require 'opsmgr/ui_helpers/config_helper'
require 'capybara'

SPEC_ROOT = File.expand_path(File.dirname(__FILE__))

Capybara.save_and_open_page_path = File.expand_path(File.join(SPEC_ROOT, '..', 'tmp'))
