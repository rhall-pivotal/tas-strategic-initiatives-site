#!/bin/bash

set -ex

script_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
workspace_dir="$( cd ${script_dir}/.. && pwd )"

docker run --privileged \
  -v ${workspace_dir}:/workspace \
  -e ACCESS_KEY_ID \
  -e SECRET_ACCESS_KEY \
  -e RUNTIME_DOT_PIVOTAL_FILE \
  -e RELENG_ENV \
  -e ENVIRONMENTS_DIR=config/environments \
  -w "/workspace/p-runtime" \
  ${DOCKER_REGISTRY}/releng:`git describe --dirty` /bin/sh -c "$*"
