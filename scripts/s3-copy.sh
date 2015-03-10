#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 S3_SOURCE S3_DEST" >&2
  echo >&2
  echo "Example: $0 s3://bucket/some/nested/key s3://bucket/another/path" >&2
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
S3_SOURCE=$1
S3_DEST=$2
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker run \
  -v $SCRIPT_DIR:/opt/workspace \
  -e "S3_KEY=${S3_KEY}" \
  -e "S3_SECRET=${S3_SECRET}" \
  ${DOCKER_REGISTRY}/releng/releng-blobstore \
  /opt/workspace/support/_s3-copy.sh $S3_SOURCE $S3_DEST
