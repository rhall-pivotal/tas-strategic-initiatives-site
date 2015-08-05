require 'net/ssh'
require 'opsmgr/environments'

module Ert
  class InternetChecker
    def initialize(environment_name:)
      @environment = Opsmgr::Environments.for(environment_name)
    end

    OPS_MANAGER_USERNAME = 'ubuntu'
    OPS_MANAGER_PASSWORD = 'tempest'

    def connection_allowed?(hostname, port)
      uri = URI.parse(environment.settings.ops_manager.url)
      result = nil
      Net::SSH.start(uri.host, OPS_MANAGER_USERNAME, password: OPS_MANAGER_PASSWORD) do |ssh|
        result = ssh_exec!(ssh, "echo QUIT | nc -w 5 #{hostname} #{port}")
      end
      result[2] == 0
    end

    def internetless?
      !connection_allowed?('www.google.com', 80) &&
        !connection_allowed?('www.google.com', 443) &&
        connection_allowed?('smtp-relay.gmail.com', 25) &&
        connection_allowed?('smtp-relay.gmail.com', 587) &&
        !connection_allowed?('github.com', 22)
    end

    private

    attr_accessor :environment

    def ssh_exec!(ssh, command)
      stdout_data = stderr_data = ''
      exit_code = exit_signal = nil
      ssh.open_channel do |channel|
        channel.exec(command) do |_ch, success|
          unless success
            abort "FAILED: couldn't execute command (ssh.channel.exec)"
          end
          channel.on_data { |_, data| stdout_data += data }

          channel.on_extended_data { |_, _, data| stderr_data += data }

          channel.on_request('exit-status') { |_, data| exit_code = data.read_long }

          channel.on_request('exit-signal') { |_, data| exit_signal = data.read_string }
        end
      end
      ssh.loop
      [stdout_data, stderr_data, exit_code, exit_signal]
    end
  end
end
