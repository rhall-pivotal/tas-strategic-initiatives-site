#!/bin/bash

if [ $# -ne 2 ]; then
  echo "Usage: $0 URL LOCAL_FILE" >&2
  echo >&2
  echo "Example: $0 http://some.url.com /path/to/file" >&2
  exit 1
fi

set -e

SCRIPT_DIR=$( cd "$( dirname $0 )" && pwd)
URL=$1
LOCAL_FILE=$2

mkdir -p $( dirname $LOCAL_FILE )
FOLDER_OF_LOCAL_FILE="$( cd "$( dirname $LOCAL_FILE )" && pwd )"
LOCAL_FILE_BASENAME=$(basename $LOCAL_FILE)
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker run \
  -v $FOLDER_OF_LOCAL_FILE:/data \
  -v $SCRIPT_DIR:/opt/workspace \
  ${DOCKER_REGISTRY}/releng/releng \
  /opt/workspace/support/_download.sh $URL /data/$LOCAL_FILE_BASENAME
