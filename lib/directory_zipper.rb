require 'zip'

class DirectoryZipper
  def initialize(target, source_directories)
    @target = target
    @source_directories = source_directories
  end

  def zip
    Zip::File.open(@target, Zip::File::CREATE) do |zipfile|
      @source_directories.each do |directory|
        working_dir = File.expand_path(File.join(directory, '..'))
        files_in_directory = File.join(directory, '**', '**')
        Dir[files_in_directory].each do |file_on_fs|
          key = file_on_fs.sub("#{working_dir}/", '')
          zipfile.add(key, file_on_fs)
        end
      end
    end
  end
end
