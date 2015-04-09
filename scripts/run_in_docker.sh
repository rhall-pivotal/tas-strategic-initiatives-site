#!/bin/bash

set -ex

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
WORKSPACE_DIR="$( cd ${SCRIPT_DIR}/../.. && pwd )"

SSH_KEY_LOCATION=${SSH_KEY_LOCATION:-/var/vcap/jobs/gocd-agent/id_rsa}

DOCKER_REGISTRY=${DOCKER_REGISTRY:-docker.vsphere.gocd.cf-app.com:5000}
DOCKER_IMAGE=${DOCKER_IMAGE:-releng/releng}

docker run \
  -v ${SSH_KEY_LOCATION}:/root/.ssh/id_rsa \
  -v ${WORKSPACE_DIR}:/workspace \
  -e AWS_ACCESS_KEY_ID \
  -e AWS_SECRET_ACCESS_KEY \
  -w "/workspace/p-runtime" \
  ${DOCKER_REGISTRY}/${DOCKER_IMAGE} $@
