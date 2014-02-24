require 'yaml'
require 'directory_zipper'
require 'open4'
require 'psych'
require 'digest'

class ProductBuilder
  def initialize(name, target_manifest, stemcell_version, working_dir, metadata)
    raise ArgumentError, 'name must be present' if name.nil? || name.empty?

    @name = name
    @target_manifest = target_manifest
    @working_dir = working_dir
    @metadata = metadata
    @stemcell_version = stemcell_version
  end

  def build
    raise "file already exists: #{pivotal_output_path}" if File.exists?(pivotal_output_path)

    create_bosh_release
    download_stemcell_if_not_local(stemcell_version)
    update_metadata
    zip
  end

  def create_bosh_release
    manifest_path = File.expand_path(File.join(bosh_release_dir, target_manifest))
    parent_path = File.expand_path(File.join(bosh_release_dir, ".."))

    raise "#{manifest_path} does not exist" unless File.exist?(manifest_path)

    tarball_filename = File.basename(project_tarball_path)
    bosh_created_tarball_path = File.join(bosh_release_dir, tarball_filename)

    if !File.exists?(bosh_created_tarball_path)
      bosh_console_output = ''
      status = Open4::popen4("cd #{parent_path} && bosh create release --with-tarball '#{manifest_path}'") do |_pid, _stdin, _stdout, stderr|
        bosh_console_output = stderr.read.strip
      end
      raise("bosh create failed: #{bosh_console_output}") unless status.success?
    end

    FileUtils.cp(bosh_created_tarball_path, project_tarball_path)
  end

  def update_metadata
    metadata['releases'] = releases_metadata
    metadata['stemcell'] = stemcell_metadata
    metadata['provides_product_versions'] = product_versions_metadata

    File.open(File.join(working_dir, 'metadata', 'cf.yml'), 'w') do |out|
      Psych.dump(metadata, out)
    end
  end

  def stemcell_filename(stemcell_version)
    "bosh-stemcell-#{stemcell_version}-vsphere-esxi-ubuntu.tgz"
  end

  def download_stemcell_if_not_local(stemcell_version)
    return if File.exists?(stemcell_path(stemcell_version))

    error = nil
    status = Open4::popen4("bosh download public stemcell #{stemcell_filename(stemcell_version)}") do |_, _, _, stderr|
      error = stderr.read.strip
    end

    raise "Could not download stemcell version #{stemcell_version}: #{error}" unless status.success?
  end

  def stemcell_path(stemcell_version)
    File.join(working_dir, 'stemcells', stemcell_filename(stemcell_version))
  end

  def pivotal_output_path
    File.join(working_dir, "p-runtime-#{metadata['product_version']}-#{name}.pivotal")
  end

  def zip
    zipper = DirectoryZipper.new(pivotal_output_path, working_dir)
    ['metadata', 'content_migrations'].map do |dir|
      zipper.add_directory(File.expand_path(File.join(working_dir, dir)))
    end
    zipper.add_file(stemcell_path(stemcell_version))
    zipper.add_file(project_tarball_path)
    zipper.zip
  end

  private

  def project_tarball_path
    tarball_filename = target_manifest.gsub(/\.yml$/, '.tgz')
    File.join(working_dir, 'releases', tarball_filename)
  end

  def stemcell_metadata
    {
      'name' => 'bosh-vsphere-esxi-ubuntu',
      'version' => stemcell_version,
      'file' => stemcell_filename(stemcell_version),
      'md5' => Digest::MD5.file(stemcell_path(stemcell_version)).hexdigest.encode('utf-8')
    }
  end

  def releases_metadata
    releases = []
    release_files = File.join(working_dir, 'releases', '*tgz')
    Dir.glob(release_files).each do |release_path|
      release_filename = File.basename(release_path)
      regex_match = release_filename.match(/\A(.*)-(\d+)\.tgz\Z/)
      next unless regex_match

      release_name = regex_match[1]
      release_version = regex_match[2]
      release_md5 = Digest::MD5.file(release_path).hexdigest

      release_metadata = {
        'file' => release_filename,
        'name' => release_name,
        'version' => release_version,
        'md5' => release_md5.encode('utf-8')
      }

      releases << release_metadata
    end

    releases
  end

  def product_versions_metadata
    provides_product_versions = (metadata['provides_product_versions'] || []).clone

    provided_version = provides_product_versions.find { |ppv| ppv['name'] == metadata['name'] }
    unless provided_version
      provides_product_versions << provided_version = { 'name' => metadata['name'] }
    end
    provided_version['version'] = metadata['product_version']

    provides_product_versions
  end

  def bosh_release_dir
    File.join(working_dir, "../cf-release/releases")
  end

  attr_reader :metadata, :name, :target_manifest, :working_dir, :stemcell_version
end
