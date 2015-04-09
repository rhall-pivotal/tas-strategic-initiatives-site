FROM ruby:2.1

RUN apt-get update
RUN apt-get -y upgrade; apt-get clean

RUN apt-get install -y s3cmd
RUN apt-get install -y aria2

ADD include/s3cfg.riak /s3cfg.riak
ADD include/s3cfg.s3 /s3cfg.s3
ADD https://github.com/rlmcpherson/s3gof3r/releases/download/v0.4.10/gof3r_0.4.10_linux_amd64.tar.gz /gof3r.tar.gz
RUN tar --strip-components 1 -C /usr/bin -xzf /gof3r.tar.gz


