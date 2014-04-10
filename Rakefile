$LOAD_PATH.unshift(File.expand_path('../lib', __FILE__))

require 'product_builder'
require 'bench_task'

namespace :product do
  desc 'Build .pivotal release candidate'
  bench_task :build, [:rc, :runtime_manifest_filename, :stemcell_version] do |t, args|
    working_dir = File.dirname(__FILE__)
    p_runtime_metadata = YAML.load_file(File.join(working_dir, 'metadata', 'cf.yml'))
    builder = ProductBuilder.new(args[:rc], args[:runtime_manifest_filename], args[:stemcell_version], working_dir, p_runtime_metadata)
    puts "Starting build for #{builder.pivotal_output_path}"
    builder.build
    puts "#{builder.pivotal_output_path} build completed"
  end
end

begin
  require 'rspec/core/rake_task'

  RSpec::Core::RakeTask.new(:spec)

  task default: :spec
rescue LoadError
end
