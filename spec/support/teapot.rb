require 'opsmgr/teapot'
require 'opsmgr/teapot/spec_helper'

RSpec::Matchers.define :be_logged_in do
  match do |username|
    req = Net::HTTP::Get.new("/teapot/login_status/#{username}")
    response = teapot_client.request(req)

    response.body == 'logged in'
  end
end

RSpec.configure do |c|
  c.include(Opsmgr::Teapot::SpecHelper, :teapot)

  c.before(:suite) do
    Opsmgr::Teapot.start(TeapotComponents.new)
  end

  c.after(:suite) do
    Opsmgr::Teapot.stop
  end

  c.around(:each, :teapot) do |example|
    WebMock.allow_net_connect!
    teapot_client.request(Net::HTTP::Get.new('/teapot/reset'))

    example.call

    WebMock.disable_net_connect!
  end
end
