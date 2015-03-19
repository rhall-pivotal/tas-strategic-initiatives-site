require 'opsmgr'
require 'opsmgr/log'

module Cmd
  class Artifacts
    include Opsmgr::Loggable

    def self.remove_cf_files(target_dir)
      Dir.chdir(target_dir) do
        globs = %w(*.pivotal *.pivotal.yml *.pivotal.md5)

        globs.each do |glob|
          Dir.glob(glob) do |path|
            log.info("Removing cf artifacts - #{path}")
            FileUtils.rm_f(path)
          end
        end
      end
    end

    def self.retrieve_cf_files(cache_key, bucket_list)
      log.info("Retrieving cf artifacts - #{cache_key}")
      bucket_list.download(key: "#{cache_key}.md5", destination_filename: "#{File.basename(cache_key)}.md5")
      bucket_list.download(key: "#{cache_key}.yml", destination_filename: "#{File.basename(cache_key)}.yml")
      bucket_list.download(key: cache_key, destination_filename: File.basename(cache_key))
      log.info('Retrieved cf artifacts')
    end
  end
end
