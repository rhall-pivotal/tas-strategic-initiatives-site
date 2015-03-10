#!/bin/bash
set -e

UPLOAD_FILE=$1
S3_PATH=$2

s3cmd -c /s3cfg.riak put $UPLOAD_FILE $S3_PATH
