#!/bin/bash

set -e

if [ -z "${RUNTIME_DOT_PIVOTAL_FILE}" ]; then
  echo 'You must set the RUNTIME_DOT_PIVOTAL_FILE environment variable!' >&2
  exit 1
fi

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"
PIVOTAL_FILE="${SCRIPTS_DIR}/../${RUNTIME_DOT_PIVOTAL_FILE}"
PIVOTAL_MD5="${SCRIPTS_DIR}/../${RUNTIME_DOT_PIVOTAL_FILE}.md5"

BUCKET_NAME=${BUCKET_NAME:-releng-products}
FOLDER_NAME=${1:-runtime}

gof3r get s3://${BUCKET_NAME}/#{FOLDER_NAME}/${RUNTIME_DOT_PIVOTAL_FILE} -p ${PIVOTAL_FILE}
gof3r get s3://${BUCKET_NAME}/#{FOLDER_NAME}/${RUNTIME_DOT_PIVOTAL_FILE}.md5 -p ${PIVOTAL_MD5}

bundle exec rake md5:validate[${PIVOTAL_FILE},${PIVOTAL_MD5}]
