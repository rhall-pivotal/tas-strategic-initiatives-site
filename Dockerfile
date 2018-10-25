FROM pivotalcfreleng/ruby

# Install ops-manifest
#   assumes ops-manifest repo was cloned into ./vendor/ops-manifest
COPY ./vendor/ /tmp/ops-manifest
WORKDIR /tmp/ops-manifest
RUN gem install ops-manifest.gem --no-ri
