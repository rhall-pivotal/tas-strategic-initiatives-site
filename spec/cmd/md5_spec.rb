require 'spec_helper'
require 'cmd/md5'
require 'tempfile'

describe Cmd::MD5 do
  around do |example|
    @source_file = Tempfile.new('source')
    @md5_file = Tempfile.new('md5-file')

    example.call

    @source_file.close!
    @md5_file.close!
  end

  context 'when the md5 matches' do
    it 'does not raise' do
      @md5_file.write('d41d8cd98f00b204e9800998ecf8427e')
      @md5_file.rewind

      expect do
        Cmd::MD5.validate_file(@source_file.path, @md5_file.path)
      end.to_not raise_error
    end
  end

  context 'when the md5 does not match' do
    it 'raises an error' do
      expect do
        Cmd::MD5.validate_file(@source_file.path, @md5_file.path)
      end.to raise_error(/MD5 mismatch/)
    end
  end
end
