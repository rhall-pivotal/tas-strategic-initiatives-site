FROM ruby:2.2

RUN apt-get update && apt-get -y upgrade && \
    apt-get install -y \
    build-essential \
    curl \
    git-core \
    ntp \
    wget \
    zip unzip \
    xvfb \
    qt5-default \
    libqt5webkit5-dev \
    libmysqlclient-dev \
    s3cmd \
    aria2 \
    && apt-get clean

RUN chmod -R a+w /usr/local/bundle
RUN chmod -R a+x /usr/local/bundle/bin

# Install gems from rubygems
ADD Gemfile.lock /tmp/Gemfile.lock
ADD Gemfile /tmp/Gemfile
ADD ci/scripts/helpers/bundle_rubygem_installer.rb /tmp/bundle_rubygem_installer.rb
RUN cd /tmp && ./bundle_rubygem_installer.rb
RUN rm -Rf /tmp/*

# We have dependency conflicts with bosh_cli
# Install it outside of our bundle context
RUN gem install bosh_cli
