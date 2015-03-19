require 'digest/md5'
require 'opsmgr/log'

module Cmd
  class MD5
    include Opsmgr::Loggable

    def self.validate_file(source_filepath, md5_filepath)
      log.info("Validating MD5 for #{source_filepath} using #{md5_filepath}")

      actual_md5 = Digest::MD5.file(source_filepath).hexdigest
      expected_md5 = File.read(md5_filepath).strip

      fail 'MD5 mismatch' unless actual_md5 == expected_md5

      log.info('Checksum validation complete. MD5s match!')
    end
  end
end
