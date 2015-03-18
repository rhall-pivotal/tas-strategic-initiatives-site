require 'rubygems'
require 'bundler/setup'

require 'krafa/client/tasks'
require 'opsmgr/tasks'

base_dir = File.dirname(__FILE__)
$LOAD_PATH << File.join(base_dir, 'lib')
load 'tasks/runtime.rake'

require 'rspec/core/rake_task'
RSpec::Core::RakeTask.new(:spec)

task default: [:spec]
