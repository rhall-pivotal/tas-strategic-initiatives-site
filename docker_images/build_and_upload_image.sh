#!/bin/bash
set -e
set -x

DOCKER_IMAGE_NAME=releng
IMAGE_FOLDER="$( cd "$( dirname ${BASH_SOURCE[0]} )"/${DOCKER_IMAGE_NAME} && pwd )"
DOCKER_REGISTRY_VSPHERE=docker.vsphere.gocd.cf-app.com:5000
DOCKER_REGISTRY_AWS=docker.gocd.cf-app.com:5000

docker build -t ${DOCKER_REGISTRY_VSPHERE}/releng/${DOCKER_IMAGE_NAME} ${IMAGE_FOLDER}
echo "Build finished for ${DOCKER_IMAGE_NAME} in vsphere"

docker push ${DOCKER_REGISTRY_VSPHERE}/releng/${DOCKER_IMAGE_NAME}
echo "Finished uploading image ${DOCKER_IMAGE_NAME} to ${DOCKER_REGISTRY_VSPHERE}"

docker build -t ${DOCKER_REGISTRY_AWS}/releng/${DOCKER_IMAGE_NAME} ${IMAGE_FOLDER}
echo "Build finished for ${DOCKER_IMAGE_NAME} in aws"

docker push ${DOCKER_REGISTRY_AWS}/releng/${DOCKER_IMAGE_NAME}
echo "Finished uploading image ${DOCKER_IMAGE_NAME} to ${DOCKER_REGISTRY_AWS}"
