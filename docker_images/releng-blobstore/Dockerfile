FROM ubuntu

RUN apt-get update
RUN apt-get -y upgrade; apt-get clean

RUN apt-get install -y s3cmd

ADD include/s3cfg.riak /s3cfg.riak
ADD include/s3cfg.s3 /s3cfg.s3

