require 'slack/post'
require 'json'
require 'yaml'
require 'typhoeus'
class ConcourseFailureReporter
  def initialize
    Slack::Post.configure(
      webhook_url: 'https://hooks.slack.com/services/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',
      username: 'FAA'
    )
  end

  def run
    (new_build_info - history).select { |b| %w(errored failed).include?(b['status']) }.each do |b|
      message = <<"MESSAGE"
Failed build #{b['id']}
Job: #{b['job_name']}
URL: #{ci_url}#{b['url']}
MESSAGE
      Slack::Post.post message
    end

    File.write('/Users/pivotal/reporter_history.json', new_build_info.to_json)
  end

  def new_build_info
    @response ||= JSON.parse(
      Typhoeus::Request.new(
        ci_url + '/api/v1/builds',
        method: :get,
        userpwd: ci_creds,
        ssl_verifypeer: false
      ).run.body
    )
  end

  def history
    JSON.parse(File.read('/Users/pivotal/reporter_history.json'))
  rescue Errno::ENOENT
    new_build_info
  end

  private

  def ci_creds
    "#{ci_info['targets']['ci']['username']}:#{ci_info['targets']['ci']['password']}"
  end

  def ci_url
    ci_info['targets']['ci']['api']
  end

  def ci_info
    @ci_info ||= YAML.load_file('/Users/pivotal/.flyrc')
  end
end
