#!/bin/bash

set -x

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

KEY_FILEPATH="${SCRIPTS_DIR}/../../cf-pivotal-artifacts-metadata/"
KEY_FILE=$(echo "${KEY_FILEPATH}/cf-pivotal.blobkey")
KEY=$(cat "${KEY_FILE}")

mkdir -p "${SCRIPTS_DIR}/../../${MATERIALS_DIR}"
export RUNTIME_PIVOTAL_FILE="${SCRIPTS_DIR}/../../${MATERIALS_DIR}/$(basename ${KEY})"

echo "Beginning download of pivotal file using key - ${KEY}"
bundle
bundle exec rake opsmgr:bucket:get[${KEY},${RUNTIME_PIVOTAL_FILE}]
echo "Finished downloading pivotal file. Saved to ${RUNTIME_PIVOTAL_FILE}"
