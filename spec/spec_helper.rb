$LOAD_PATH << File.expand_path('../lib', File.dirname(__FILE__))

require 'rspec'

def fixture_dir
  File.join(File.dirname(__FILE__), 'fixtures')
end

def temp_dir
  File.join(File.dirname(__FILE__), 'tmp')
end

def temp_project_dir
  File.join(temp_dir, 'project_dir')
end
