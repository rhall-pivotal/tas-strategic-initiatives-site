#!/bin/bash

set -e

RUNTIME_DIR=$( cd "$( dirname $0 )"/../.. && pwd)
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker run \
  -v ${RUNTIME_DIR}:/opt/workspace \
  -v /var/vcap/jobs/gocd-agent/id_rsa:/id_rsa \
  -e RELENG_ENV \
  -e "GIT_SSH=/opt/workspace/scripts/util/docker_ssh" \
  ${DOCKER_REGISTRY}/releng/releng \
  /opt/workspace/scripts/runtime/aws_post_install_test.sh
