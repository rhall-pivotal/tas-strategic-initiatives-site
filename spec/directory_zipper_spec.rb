require 'spec_helper'
require 'directory_zipper'

describe DirectoryZipper do
  let(:directories_to_copy) {
    ['content_migrations', 'metadata', 'stemcells'].map do |directory|
      File.join(temp_dir, 'project_dir', directory)
    end
  }
  let(:target_zip) { File.join(temp_dir, 'test.zip') }

  let(:good_file_names) do
    %w(
        metadata/cf.yml
        content_migrations/migration1.yml
        stemcells/bosh-stemcell-1111111-vsphere-esxi-ubuntu.tgz
      )
  end

  before do
    FileUtils.cp_r(File.join(fixture_dir, '.'), temp_dir)
  end

  after do
    FileUtils.rm_rf(temp_dir)
  end

  it "includes all of the passed in directories and their contents, and nothing else" do
    directory_zipper = DirectoryZipper.new(target_zip, directories_to_copy)
    directory_zipper.zip

    zip_contents = Zip::File.open(target_zip) do |zipfile|
      zipfile.map { |entry| entry.name }
    end

    expect(zip_contents).to match_array(good_file_names)
  end
end
