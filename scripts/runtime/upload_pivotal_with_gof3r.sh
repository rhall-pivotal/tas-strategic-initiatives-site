#!/bin/bash

set -ex

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

PIVOTAL_DIR="${SCRIPTS_DIR}/.."
PIVOTAL_FILE=$(echo ${PIVOTAL_DIR}/*.pivotal)
PIVOTAL_BASENAME=$(basename ${PIVOTAL_FILE})

BUCKET_NAME=${BUCKET_NAME:-releng-products}
FOLDER_NAME=${1:-untested}

gof3r cp ${PIVOTAL_FILE}.md5 s3://${BUCKET_NAME}/${FOLDER_NAME}/${PIVOTAL_BASENAME}.md5
gof3r cp ${PIVOTAL_FILE}.yml s3://${BUCKET_NAME}/${FOLDER_NAME}/${PIVOTAL_BASENAME}.yml
gof3r cp ${PIVOTAL_FILE} s3://${BUCKET_NAME}/${FOLDER_NAME}/${PIVOTAL_BASENAME}

echo "${FOLDER_NAME}/${PIVOTAL_BASENAME}.md5" > "${PIVOTAL_DIR}/cf-pivotal.md5.blobkey"
echo "${FOLDER_NAME}/${PIVOTAL_BASENAME}.yml" > "${PIVOTAL_DIR}/cf-pivotal.yml.blobkey"
echo "${FOLDER_NAME}/${PIVOTAL_BASENAME}" > "${PIVOTAL_DIR}/cf-pivotal.blobkey"
