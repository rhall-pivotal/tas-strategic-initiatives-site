source 'https://rubygems.org'

gem 'krafa-client', '0.0.11', git: 'git@github.com:pivotal-cf-experimental/krafa.git'
gem 'vara',     '0.12.0',   require: false, git: 'git@github.com:pivotal-cf/vara.git'
gem 'bosh_cli', '1.2818.0', require: false # a dependency of vara but version should match bosh stemcell: metadata_parts/binaries.yml

gem 'opsmgr', git: 'git@github.com:pivotal-cf/opsmgr', branch: 'features/add-stemcell-import-97032016'
gem 'vsphere_clients', git: 'git@github.com:pivotal-cf-experimental/vsphere_clients'
gem 'recursive-open-struct'

group :development, :test do
  gem 'rspec'
  gem 'webmock'
  gem 'byebug'
  gem 'rubocop'
end
