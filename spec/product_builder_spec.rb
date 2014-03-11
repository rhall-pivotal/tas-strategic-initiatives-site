require 'spec_helper'
require 'product_builder'
require 'open4'

describe ProductBuilder do
  around do |test|
    FileUtils.rm_rf(temp_dir)
    FileUtils.cp_r(File.join(fixture_dir, '.'), temp_dir)
    test.run
    FileUtils.rm_rf(temp_dir)
  end

  let(:target_manifest) { 'cf-9999.yml' }
  let(:seed_version_number) { '1111111' }
  let(:build_name) { "rc#{rand(10)}" }
  let(:metadata) { YAML.load_file(File.join(temp_project_dir, 'metadata', 'cf.yml')) }

  describe '.new' do
    it 'raises if given a nil name' do
      expect { ProductBuilder.new(nil, target_manifest, seed_version_number, temp_project_dir, {}) }.to raise_error(/name must be present/)
    end

    it 'raises if given an empty name' do
      expect { ProductBuilder.new('', target_manifest, seed_version_number, temp_project_dir, {}) }.to raise_error(/name must be present/)
    end
  end

  let(:compiled_package_file) { File.join(temp_project_dir, 'compiled_packages', compiled_package) }

  describe '#build' do
    let(:builder) { ProductBuilder.new(build_name, target_manifest, seed_version_number, temp_project_dir, metadata) }

    describe 'output' do
      it 'should create a zipfile' do
        expect {
          builder.build
        }.to change { File.exists?(builder.pivotal_output_path) }.from(false).to(true)
      end

      describe 'zipfile contents' do
        let(:pivotal_filepath) { builder.pivotal_output_path }
        let(:pivotal_filename) { File.basename(pivotal_filepath) }
        let(:stemcell_filename) { builder.stemcell_filename(seed_version_number) }
        let(:unzip_dir) { File.expand_path(File.join(temp_dir, 'unzip_here')) }

        it 'should explode to a properly formed directory structure' do
          builder.build
          FileUtils.mkdir_p(unzip_dir)
          `unzip #{pivotal_filepath} -d #{unzip_dir}`

          Dir.chdir(unzip_dir) do
            directory_contents = Dir.glob("**/**")
            expect(directory_contents).to match_array(%W(
              content_migrations
              content_migrations/migration1.yml
              metadata
              metadata/cf.yml
              releases
              releases/cf-9999.tgz
              stemcells
              stemcells/#{stemcell_filename}
            ))
          end
        end

        context 'with compiled packages' do
          before do
            FileUtils.touch(compiled_package_file)
          end

          let(:compiled_package) { "cf-9999-bosh-vsphere-esxi-ubuntu-#{seed_version_number}.tgz" }

          it 'should explode to a properly formed directory structure including the compiled packages' do
            builder.build
            FileUtils.mkdir_p(unzip_dir)
            `unzip #{pivotal_filepath} -d #{unzip_dir}`

            Dir.chdir(unzip_dir) do
              directory_contents = Dir.glob("**/**")
              expect(directory_contents).to match_array(%W(
                compiled_packages
                compiled_packages/#{compiled_package}
                content_migrations
                content_migrations/migration1.yml
                metadata
                metadata/cf.yml
                releases
                releases/cf-9999.tgz
                stemcells
                stemcells/#{stemcell_filename}
              ))
            end
          end
        end
      end
    end

    describe 'validations' do
      let(:build_name) { 'preexisting' }

      it 'should raise if target .pivotal already exists' do
        expect {
          builder.build
        }.to raise_error(/file already exists/)
      end
    end

    context 'with a -dev.yml manifest file' do
      let(:target_manifest) { 'cf-9999.9999-dev.yml' }
      it 'should look in the dev_releases/ dir and not the releases/ dir for the .yml and create the .tgz there' do
        expect {
          builder.build
        }.to change { File.exists?(builder.pivotal_output_path) }.from(false).to(true)
        expect(File.exist?(File.join(temp_dir, 'cf-release', 'dev_releases', 'cf-9999.9999-dev.tgz'))).to be_true
      end
    end
  end

  describe '#create_bosh_release' do
    let(:builder) { ProductBuilder.new(build_name, target_manifest, seed_version_number, temp_project_dir, metadata) }

    context 'if the target manifest is missing' do
      let(:target_manifest) { "manifest_that_doesnt_exist_version_#{rand(10)}.yml" }

      it 'should raise' do
        expect { builder.create_bosh_release }.to raise_error(/#{target_manifest} does not exist/)
      end
    end

    context 'if the target cf manifest is present' do
      let(:tarball_filename) { target_manifest.gsub(/\.yml$/, '.tgz') }
      let(:project_tarball_path) { File.join(temp_project_dir, 'releases', tarball_filename) }
      let(:bosh_create_tarball_path) { File.join(temp_project_dir, '../cf-release/releases', tarball_filename) }
      let(:tarball_contents) { "I am some random bits in a tarball #{rand(10)}" }
      let(:tarball_md5) { Digest::MD5.hexdigest(tarball_contents) }

      context 'if the target mysql tarball is present' do
        before do
          File.write(bosh_create_tarball_path, tarball_contents)
          Open4.stub(:popen4).and_return(double 'status', success?: true)
          builder.create_bosh_release
        end

        it 'should not shell out to bosh' do
          expect(Open4).not_to have_received(:popen4)
        end

        it 'should copy the release tarball to local releases/' do
          expect(File.exist?(project_tarball_path)).to be_true
          expect(FileUtils.cmp(project_tarball_path, bosh_create_tarball_path)).to be_true
        end
      end

      context 'if the target mysql tarball is missing' do
        it 'should execute a bosh create' do
          allow(Open4).to receive(:popen4).and_call_original

          builder.create_bosh_release
          cf_release_dir = File.expand_path(File.join(temp_dir, 'cf-release'))
          expect(Open4).to have_received(:popen4).with("cd #{cf_release_dir} && bosh create release --with-tarball '#{File.join(cf_release_dir, 'releases', target_manifest)}'")
        end

        context 'bosh succeeds' do
          it 'should copy the release tarball to local releases/' do
            builder.create_bosh_release

            expect(File.exist?(project_tarball_path)).to eq(true)
            expect(FileUtils.cmp(project_tarball_path, bosh_create_tarball_path)).to eq(true)
          end
        end

        context 'bosh fails' do
          before do
            allow(Open4).to receive(:popen4) do |command, &blk|
              stderr = double(read: double(strip: "bosh create failed"))
              blk.call(nil, nil, nil, stderr)

              double(success?: false)
            end
          end

          it 'should raise' do
            expect { builder.create_bosh_release }.to raise_error(/bosh create failed/)
          end
        end
      end
    end
  end

  describe '#update_metadata' do
    let(:original_metadata) do
      {
        'foo' => 'bar',
        'name' => 'p-foo',
        'product_version' => '1.2.3.4',
        'provides_product_versions' => [
          {'name' => 'whatever', 'version' => '0.0.0'},
          {'name' => 'p-foo', 'version' => '1.2.1.1'},
          {'name' => 'p-old-foo', 'version' => '10.0.2'},
        ],
        'compiled_package' => {
          'name' => 'fake-compiled-package'
        }
      }
    end
    let(:target_manifest) { 'cf-123.yml' }
    let(:builder) { ProductBuilder.new(build_name, target_manifest, seed_version_number, temp_project_dir, original_metadata) }
    let(:updated_metadata) { YAML.load_file(File.join(temp_project_dir, 'metadata', 'cf.yml')) }

    it 'should copy keys from the original metadata' do
      builder.update_metadata
      expect(updated_metadata['foo']).to eq('bar')
    end

    it 'should update the releases array with target release' do
      builder.update_metadata
      expected_release_metadata = {
        'name' => 'cf',
        'file' => 'cf-123.tgz',
        'version' => '123',
        'md5' => 'c3d48e9a5e5bd66dc09b3009cead694c'
      }
      expect(updated_metadata['releases']).to eq([expected_release_metadata])
    end

    it 'should update the stemcell info with the new version, filename, and md5' do
      builder.update_metadata
      expected_stemcell_metadata = {
        'name' => 'bosh-vsphere-esxi-ubuntu',
        'version' => seed_version_number,
        'file' => "bosh-stemcell-#{seed_version_number}-vsphere-esxi-ubuntu.tgz",
        'md5' => 'd41d8cd98f00b204e9800998ecf8427e'
      }
      expect(updated_metadata['stemcell']).to eq(expected_stemcell_metadata)
    end

    context 'with a compiled_packages tarball' do
      before do
        puts compiled_package_file
        File.write(compiled_package_file, 'I am a compiled package, howdy yall')
      end

      let(:compiled_package) { 'cf-123-bosh-vsphere-esxi-ubuntu-1111111.tgz' }

      it 'should re-compute the compiled_packages info' do
        builder.update_metadata
        compiled_package_metadata = {
          'name' => 'cf',
          'version' => '123',
          'file' => compiled_package,
          'md5' => '91f491bdb83c77ee16abd69323f42ad7'
        }
        expect(updated_metadata['compiled_package']).to eq(compiled_package_metadata)
      end
    end

    context 'without a compiled_packages tarball' do
      it 'should remove any compiled_packages info from the metadata' do
        builder.update_metadata
        expect(updated_metadata['compiled_package']).to be_nil
      end
    end

    describe 'provides_product_versions' do
      it 'synchronizes product_version with provides_product_versions[n][version] when the names match' do
        builder.update_metadata
        expected_provides_product_versions_metadata = [
          {'name' => 'whatever', 'version' => '0.0.0'},
          {'name' => 'p-foo', 'version' => '1.2.3.4'},
          {'name' => 'p-old-foo', 'version' => '10.0.2'}
        ]
        expect(updated_metadata['provides_product_versions']).to eq(expected_provides_product_versions_metadata)
      end

      context 'when there was no expected_provides_versions key in the original metadata' do
        let(:original_metadata) do
          {
            'name' => 'p-foo',
            'product_version' => '1.2.3.4'
          }
        end

        it 'adds the name and product_version to the expected_provides_versions key' do
          builder.update_metadata
          expected_provides_product_versions_metadata = [
            {'name' => 'p-foo', 'version' => '1.2.3.4'}
          ]
          expect(updated_metadata['provides_product_versions']).to eq(expected_provides_product_versions_metadata)
        end
      end
    end
  end

  describe '#stemcell_filename' do
    it 'should return the stemcell filename from a version' do
      version = 45
      filename = "bosh-stemcell-#{version}-vsphere-esxi-ubuntu.tgz"
      pb = ProductBuilder.new('rc-test', target_manifest, seed_version_number, temp_project_dir, {})
      expect(pb.stemcell_filename(version)).to eq(filename)
    end
  end

  describe '#download_stemcell_if_not_local' do
    context 'when the stemcell exists locally' do
      it 'uses the local stemcell' do
        allow(Open4).to receive(:popen4)
        pb = ProductBuilder.new('rc-test', target_manifest, seed_version_number, temp_project_dir, {})
        pb.download_stemcell_if_not_local seed_version_number
        expect(Open4).not_to have_received(:popen4)
      end
    end

    context 'when the stemcell does not exist locally' do
      let(:seed_version_number) { '99999' }
      let(:stemcell_folder) { File.expand_path(File.join(temp_project_dir, 'stemcells')) }

      it 'downloads the file from bosh_artifacts' do
        allow(Open4).to receive(:popen4).and_return(double('status_double', success?: true))

        pb = ProductBuilder.new('rc-test', target_manifest, seed_version_number, temp_project_dir, {})
        pb.download_stemcell_if_not_local seed_version_number
        expect(Open4).to have_received(:popen4).with("mkdir -p #{stemcell_folder} && cd #{stemcell_folder} && bosh download public stemcell bosh-stemcell-#{seed_version_number}-vsphere-esxi-ubuntu.tgz")
      end

      context 'when downloading fails' do
        it 'raises' do
          allow(Open4).to receive(:popen4) do |command, &blk|
            stderr = double(read: double(strip: "external error message"))
            blk.call(nil, nil, nil, stderr)

            double(success?: false)
          end

          pb = ProductBuilder.new('rc-test', target_manifest, seed_version_number, temp_project_dir, {})
          expect {
            pb.download_stemcell_if_not_local seed_version_number
          }.to raise_error(/Could not download.*external error message/)
        end
      end
    end
  end

  describe '#pivotal_output_path' do
    let(:release_name) { build_name }
    let(:product_version) { "#{rand(3)}.2.1.0" }
    let(:metadata) { {'product_version' => product_version} }
    let(:product_builder) { ProductBuilder.new(release_name, target_manifest, seed_version_number, temp_project_dir, metadata) }

    it 'should compose the path from workdir, product version, and given name' do
      expected_path = File.join(temp_project_dir, "cf-#{product_version}-#{release_name}.pivotal")
      expect(product_builder.pivotal_output_path).to eq(expected_path)
    end
  end
end
