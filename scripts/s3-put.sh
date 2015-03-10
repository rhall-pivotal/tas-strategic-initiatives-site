#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 UPLOAD_FILE S3_PATH" >&2
  echo >&2
  echo "Example: $0 /path/to/file s3://bucket/some/nested/key" >&2
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
UPLOAD_FILE=$1
S3_PATH=$2

FOLDER_OF_UPLOAD_FILE="$( cd "$( dirname $UPLOAD_FILE )" && pwd )"
UPLOAD_FILE_BASENAME=$(basename $UPLOAD_FILE)
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker run \
  -v $FOLDER_OF_UPLOAD_FILE:/data \
  -v $SCRIPT_DIR:/opt/workspace \
  -e "S3_KEY=${S3_KEY}" \
  -e "S3_SECRET=${S3_SECRET}" \
  ${DOCKER_REGISTRY}/releng/releng-blobstore \
  /opt/workspace/support/_s3-put.sh /data/$UPLOAD_FILE_BASENAME $S3_PATH
