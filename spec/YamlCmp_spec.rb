require 'rspec'
require 'psych'

describe 'CfYmlChanges' do
  before do
    `cp metadata/cf.yml /tmp && git checkout metadata/cf.yml`
  end
  after do
    `cp /tmp/cf.yml metadata/`
  end

  xit 'this test will fail and the output in RubyMine will help you verify that the changes are appropriate' do
    modified = Psych.load_file('/tmp/cf.yml')
    original = Psych.load_file('metadata/cf.yml')
    expect(original).to eq(modified)
  end
end
