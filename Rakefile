require 'rubygems'
require 'bundler/setup'

require 'krafa/client/tasks'
require 'opsmgr/tasks'

base_dir = File.expand_path(File.dirname(__FILE__))
$LOAD_PATH << File.join(base_dir, 'lib')

Dir.glob(File.join(base_dir, 'lib', 'tasks', '*.rake')).each do |tasks_file|
  load(tasks_file)
end

require 'rspec/core/rake_task'
RSpec::Core::RakeTask.new(:spec)

task default: [:spec]
