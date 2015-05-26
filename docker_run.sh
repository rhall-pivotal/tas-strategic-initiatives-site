#!/bin/bash

set -ex

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
workspace_dir="$( cd ${script_dir}/.. && pwd )"

docker run --privileged \
  -v ${workspace_dir}:/workspace \
  -e ACCESS_KEY_ID \
  -e SECRET_ACCESS_KEY \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -e RIAK_BUCKET \
  -e S3_BUCKET \
  -e RUNTIME_DOT_PIVOTAL_FILE \
  -e RELENG_ENV \
  -e VARA_EXTRA_FLAGS \
  -w "/workspace/p-runtime" \
  ${DOCKER_REGISTRY}/releng:`git describe --dirty` /bin/sh -c "$*"
