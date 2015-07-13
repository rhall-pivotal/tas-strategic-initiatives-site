require 'yaml'
require 'active_support/core_ext/hash/deep_merge'

files = ['aws-clean.yml', 'vsphere-clean.yml', 'vsphere-internetless.yml', 'aws-upgrade.yml', 'vsphere-upgrade.yml']

result = YAML.load_file('ert-1.5.yml')
files.each do |f|
  result['jobs'] += YAML.load_file(f)['jobs']
end

output = YAML.dump(result)
File.open('ert-1.5-full-pipeline.yml', 'w') do |f|
  f.write(output)
end

puts 'ERT-1.5 pipeline config written to ert-1.5-full-pipeline.yml'
