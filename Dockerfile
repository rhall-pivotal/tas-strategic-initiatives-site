FROM pivotalcfreleng/ruby
MAINTAINER https://github.com/pivotal-cf/p-runtime

COPY ./vendor/ /tmp/ops-manifest
WORKDIR /tmp/ops-manifest
RUN gem install ops-manifest.gem --no-ri
