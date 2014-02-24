require 'zip'

class DirectoryZipper
  def initialize(target, common_source_directory)
    @target = target
    @common_source_directory = common_source_directory
    @source_directories = []
    @extra_includes = []
  end

  def add_directory(directory_path)
    check_common_directory(directory_path)
    @source_directories << directory_path
  end

  def add_file(file_path)
    check_common_directory(file_path)
    @extra_includes << file_path
  end

  def zip
    Zip::File.open(@target, Zip::File::CREATE) do |zipfile|
      @source_directories.each do |directory|
        zip_item(zipfile, directory)

        files_in_directory = File.join(directory, '**', '**')
        Dir[files_in_directory].each do |file_path|
          zip_item(zipfile, file_path)
        end
      end

      @extra_includes.each do |extra_include_path|
        zip_item(zipfile, extra_include_path)
      end
    end
  end

  private

  def zip_item(open_zipfile, path)
    relative_path = key_for_path(path)
    return if relative_path.empty?

    relative_directory = File.dirname(relative_path)
    if relative_directory != '.' && open_zipfile.find_entry(relative_directory).nil?
      open_zipfile.add(relative_directory, File.join(@common_source_directory, relative_directory))
    end
    open_zipfile.add(relative_path, path) if open_zipfile.find_entry(relative_path).nil?
  end

  def key_for_path(path)
    key = path.sub(@common_source_directory, '')
    key.gsub(/^\//, '')
  end

  def check_common_directory(path)
    raise("Must add files and directories under the common source directory") unless path.start_with?(@common_source_directory)
  end
end
