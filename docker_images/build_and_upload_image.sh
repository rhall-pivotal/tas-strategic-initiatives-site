#!/bin/bash
set -e

DOCKER_IMAGE_NAME=releng-blobstore
IMAGE_FOLDER="$( cd "$( dirname ${BASH_SOURCE[0]} )"/${DOCKER_IMAGE_NAME} && pwd )"
DOCKER_REGISTRY=docker.vsphere.gocd.cf-app.com:5000

docker build -t ${DOCKER_REGISTRY}/releng/${DOCKER_IMAGE_NAME} ${IMAGE_FOLDER}
echo "Build finished for ${DOCKER_IMAGE_NAME}"

docker push ${DOCKER_REGISTRY}/releng/${DOCKER_IMAGE_NAME}
echo "Finished uploading image ${DOCKER_IMAGE_NAME} to ${DOCKER_REGISTRY}"
