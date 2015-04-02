require 'openssl'
require 'securerandom'

module Tools
  class SelfSignedRsaCertificate
    # rubocop:disable Metrics/AbcSize
    # rubocop:disable Metrics/MethodLength
    def self.generate(wildcard_domains)
      private_key = OpenSSL::PKey::RSA.new(2048)

      csr = OpenSSL::X509::Request.new
      csr.version = 0
      csr.subject = OpenSSL::X509::Name.parse("/C=US/O=Pivotal/CN=#{wildcard_domains.first}")
      csr.public_key = private_key.public_key
      csr.sign(private_key, OpenSSL::Digest::SHA1.new)

      cert = OpenSSL::X509::Certificate.new
      cert.serial = SecureRandom.random_number(2**(8 * 20))
      cert.version = 2
      cert.not_before = Time.now - (60 * 60 * 24) # yesterday
      cert.not_after = cert.not_before + (2 * 365 * 24 * 60 * 60) # 2 years
      cert.subject = csr.subject
      cert.public_key = csr.public_key
      cert.issuer = csr.subject # self-signed

      ef = OpenSSL::X509::ExtensionFactory.new
      ef.subject_certificate = cert
      ef.issuer_certificate = cert

      dns_string = wildcard_domains.map { |d| "DNS:#{d}" }.join(', ')
      ext = ef.create_extension('subjectAltName', dns_string, false)

      cert.add_extension(ext)
      cert.sign(private_key, OpenSSL::Digest::SHA1.new)

      new(private_key.to_pem, cert.to_pem)
    end
    # rubocop:enable Metrics/MethodLength
    # rubocop:enable Metrics/AbcSize

    def initialize(private_key_pem, cert_pem)
      @private_key_pem = private_key_pem
      @cert_pem = cert_pem
    end

    attr_reader :private_key_pem, :cert_pem
  end
end
