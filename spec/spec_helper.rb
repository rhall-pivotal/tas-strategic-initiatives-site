SPEC_ROOT = File.expand_path(File.dirname(__FILE__))

$LOAD_PATH << File.expand_path('../lib', SPEC_ROOT)

require 'rspec'
# require 'opsmgr/log'

# Opsmgr::Log.test_mode!

Dir['./spec/support/**/*.rb'].each { |f| require f }
