#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 S3_PATH LOCAL_FILE" >&2
  echo >&2
  echo "Example: $0 s3://bucket/some/nested/key /path/to/file" >&2
  exit 1
fi

if [ -z "$S3_KEY" ]; then
  echo "S3_KEY environment variable not set; aborting" >&2
  exit 1
fi

if [ -z "$S3_SECRET" ]; then
  echo "S3_SECRET environment variable not set; aborting" >&2
  exit 1
fi

set -e

SCRIPT_DIR=$( cd "$( dirname $0 )" && pwd)
S3_PATH=$1
LOCAL_FILE=$2

mkdir -p $( dirname $LOCAL_FILE )
FOLDER_OF_LOCAL_FILE="$( cd "$( dirname $LOCAL_FILE )" && pwd )"
LOCAL_FILE_BASENAME=$(basename $LOCAL_FILE)
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker run \
  -v $FOLDER_OF_LOCAL_FILE:/data \
  -v $SCRIPT_DIR:/opt/workspace \
  -e "S3_KEY=${S3_KEY}" \
  -e "S3_SECRET=${S3_SECRET}" \
  ${DOCKER_REGISTRY}/releng/releng-blobstore \
  /opt/workspace/support/_s3-get.sh $S3_PATH /data/$LOCAL_FILE_BASENAME
