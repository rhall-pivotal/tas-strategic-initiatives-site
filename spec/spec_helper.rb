SPEC_ROOT = File.expand_path(File.dirname(__FILE__))

$LOAD_PATH << File.expand_path('../lib', SPEC_ROOT)

require 'rspec'
require 'timeout'
require 'webmock/rspec'
require 'opsmgr/log'

require 'capybara'

Opsmgr::Log.test_mode!

Dir['./spec/support/*.rb'].each { |f| require f }

RSpec.configure do |config|
  Capybara.save_and_open_page_path = File.expand_path(File.join(SPEC_ROOT, '..', 'tmp'))

  config.after(:suite) { WebMock.disable! } # for codeclimate coverage reporting
end
