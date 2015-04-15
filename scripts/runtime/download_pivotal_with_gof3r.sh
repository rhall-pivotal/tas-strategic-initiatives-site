#!/bin/bash

set -e

if [ -z "${RUNTIME_DOT_PIVOTAL_FILE}" ]; then
  echo 'You must set the RUNTIME_DOT_PIVOTAL_FILE environment variable!' >&2
  exit 1
fi

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"
PIVOTAL_FILE="${SCRIPTS_DIR}/../${RUNTIME_DOT_PIVOTAL_FILE}"

BUCKET_NAME="releng-products"
FOLDER_NAME="runtime"

# NOTE: gof3r automatically checks MD5 by checking for a .md5 file with the same name
gof3r cp s3://${BUCKET_NAME}/${FOLDER_NAME}/${RUNTIME_DOT_PIVOTAL_FILE} ${PIVOTAL_FILE}
