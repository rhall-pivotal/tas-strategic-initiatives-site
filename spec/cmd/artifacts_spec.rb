require 'spec_helper'
require 'cmd/artifacts'
require 'bucket_brigade/bucket_list'
require 'tempfile'

describe Cmd::Artifacts do
  context '.remove_cf_files' do
    around do |example|
      Dir.mktmpdir('artifacts') do |dir|
        @target_dir = dir
        @piv = File.join(dir, 'cf.pivotal')
        @md5 = File.join(dir, 'cf.pivotal.md5')
        @yml = File.join(dir, 'cf.pivotal.yml')

        FileUtils.touch([@piv, @md5, @yml])
        example.call
      end
    end

    let(:cf_file_glob) { '*{.pivotal,.pivotal.yml,.pivotal.md5}' }

    context 'and trying to remove them' do
      it 'removes the files' do
        expect(File).to exist(@piv)
        expect(File).to exist(@md5)
        expect(File).to exist(@yml)

        Cmd::Artifacts.remove_cf_files(@target_dir)

        expect(File).not_to exist(@piv)
        expect(File).not_to exist(@md5)
        expect(File).not_to exist(@yml)
      end
    end
  end

  context '.retrieve_cf_files' do
    it 'fetches the correct files using BucketBrigade' do
      cache_key = 'folder/subfolder/cf-123.pivotal'
      fake_storage = instance_double(BucketBrigade::BucketList)

      expect(fake_storage).to receive(:download).with(key: cache_key, destination_filename: 'cf-123.pivotal')
      expect(fake_storage).to receive(:download).with(key: "#{cache_key}.md5", destination_filename: 'cf-123.pivotal.md5')
      expect(fake_storage).to receive(:download).with(key: "#{cache_key}.yml", destination_filename: 'cf-123.pivotal.yml')

      Cmd::Artifacts.retrieve_cf_files(cache_key, fake_storage)
    end
  end
end
