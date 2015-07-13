#!/usr/bin/env ruby

require 'bundler'

lockfile = File.read('Gemfile.lock')
parser = Bundler::LockfileParser.new(lockfile)

rubygems = parser.specs.select { |spec| spec.source.class == Bundler::Source::Rubygems }
rubygems.each do |spec|
  gem = "#{spec.name} -v #{spec.version.version}"
  p "Installing #{gem}"
  `gem install #{gem}`
end
