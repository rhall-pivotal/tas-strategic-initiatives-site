require 'spec_helper'
require 'dot_pivotal'

describe DotPivotal do
  let(:base_dir) { File.expand_path(File.join(SPEC_ROOT, '..')) }

  subject(:dot_pivotal) { described_class.new(base_dir) }

  describe '#base_dir' do
    it 'returns the directory passed in' do
      expect(dot_pivotal.base_dir).to eq(base_dir)
    end
  end

  describe '#metadata_file' do
    it 'defaults to nil' do
      expect(dot_pivotal.metadata_file).to be_nil
    end

    context 'after #build_metadata' do
      let(:metadata_file) { 'path/to/metadata.yml'}

      before do
        allow(Vara::ProductMetadataBuilder).to receive(:build).and_return(metadata_file)
        dot_pivotal.build_metadata
      end

      it 'returns the location of the generated metadata' do
        expect(dot_pivotal.metadata_file).to eq(metadata_file)
      end
    end
  end

  describe '#build_metadata' do
    it 'uses Vara::ProductMetadataBuilder to generate the metadata file' do
      expect(Vara::ProductMetadataBuilder).to receive(:build).with(base_dir)

      dot_pivotal.build_metadata
    end

    it 'sets #metadata_file from Vara::ProductMetadataBuilder' do
      allow(Vara::ProductMetadataBuilder).to receive(:build).and_return('path/to/metadata.yml')

      expect { dot_pivotal.build_metadata }.to change { dot_pivotal.metadata_file }.from(nil).to('path/to/metadata.yml')
    end
  end

  context 'when the metadata file exists' do
    let(:product_metadata) { double(Vara::ProductMetadata) }

    before do
      metadata_file = 'path/to/metadata.yml'
      allow(Vara::ProductMetadataBuilder).to receive(:build).and_return(metadata_file)
      dot_pivotal.build_metadata


      allow(Vara::ProductMetadata).to receive(:from_file).with(metadata_file).and_return(product_metadata)
    end

    describe '#stemcell_version' do
      it 'uses the stemcell information from the product metadata' do
        stemcell_metadata = double(Vara::StemcellMetadata)
        allow(stemcell_metadata).to receive(:version).and_return('8675309')
        allow(product_metadata).to receive(:stemcell_metadata).and_return(stemcell_metadata)

        expect(dot_pivotal.stemcell_version).to eq('8675309')
      end
    end

    describe '#releases' do
      it 'uses the releases information from the product metadata' do
        allow(product_metadata).to receive(:releases_metadata).and_return(
                                     [
                                       {
                                         'file' => 'my-release-1024.tgz',
                                         'name' => 'my-release',
                                         'version' => '1024',
                                         'md5' => ' deadbeefdeadbeefdeadbeefdeadbeef',
                                         'url' => 'https://example.com/my-release-1024.tgz'
                                       },
                                       {
                                         'file' => 'your-release-2048.tgz',
                                         'name' => 'your-release',
                                         'version' => '2048',
                                         'md5' => 'beefdeadbeefdeadbeefdeadbeefdead',
                                         'url' => 'https://example.com/your-release-2048.tgz'
                                       },
                                     ]
                                   )

        expect(dot_pivotal.releases).to eq([
                                             {'name' => 'my-release', 'version' => '1024'},
                                             {'name' => 'your-release', 'version' => '2048'},
                                           ])
      end
    end

    describe '#metadata_ref' do
      it 'fetches the current git ref of the product metadata repo' do
        expected_sha = 'feedbad'

        expect(dot_pivotal).to receive(:`).with("cd #{base_dir} && git log -1 --format=format:%h").and_return(expected_sha)

        expect(dot_pivotal.metadata_ref).to eq(expected_sha)
      end
    end

    describe '#build_identifier' do
      context 'when NO build identifier is provided' do
        it 'sets the build identifier to "local"' do
          expect(dot_pivotal.build_identifier).to eq('local')
        end
      end

      context 'when a build identifier is provided' do
        subject(:dot_pivotal) { described_class.new(base_dir, 'THE_BEST_BUILD') }


        it 'knows the build CI server build identifier' do
          expect(dot_pivotal.build_identifier).to eq('THE_BEST_BUILD')
        end
      end
    end
  end
end
