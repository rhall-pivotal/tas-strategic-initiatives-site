FROM ruby:2.1

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
    && apt-get clean

RUN apt-get install -y s3cmd
RUN apt-get install -y aria2

ADD include/s3cfg.riak /s3cfg.riak
ADD include/s3cfg.s3 /s3cfg.s3
ADD include/gof3r.tar.gz /gof3r
RUN mv /gof3r/*/gof3r /usr/bin

# Create .ssh directory so we can add keys etc
RUN mkdir -p /root/.ssh

# Ignore ssh fingerprints
RUN echo "Host * \n\tStrictHostKeyChecking no \n\tUserKnownHostsFile=/dev/null" >> /root/.ssh/config
ADD include/id_rsa /root/.ssh/id_rsa

ADD Gemfile /Gemfile
ADD Gemfile.lock /Gemfile.lock

RUN ssh-agent bundle

#Don't publish ssh key
RUN rm /root/.ssh/id_rsa
