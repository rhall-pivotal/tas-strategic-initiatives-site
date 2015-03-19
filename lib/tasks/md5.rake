namespace :md5 do
  desc "Validate a file's md5 against a given md5"
  task :validate, [:file, :md5_file] do |_, args|
    require 'cmd/md5'

    Cmd::MD5.validate_file(File.expand_path(args.file), File.expand_path(args.md5_file))
  end
end
