#!/bin/bash
set -e

S3_PATH=$1
LOCAL_FILE=$2

s3cmd -c /s3cfg.riak get --force $S3_PATH $LOCAL_FILE
