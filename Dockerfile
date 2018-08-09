FROM pivotalcfreleng/ruby
MAINTAINER https://github.com/pivotal-cf/p-runtime

# Install ops-manifest
#   assumes ops-manifest repo was cloned into ./vendor/ops-manifest
COPY ./vendor/ops-manifest /tmp/ops-manifest
WORKDIR /tmp/ops-manifest
RUN bundle install && \
  rm -rf *.gem && \
  gem build ops-manifest.gemspec && \
  gem install ops-manifest-*.gem
