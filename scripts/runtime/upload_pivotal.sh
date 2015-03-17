#!/bin/bash

set -ex

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

PIVOTAL_DIR="${SCRIPTS_DIR}/../../p-runtime"
PIVOTAL_FILE=$(echo ${PIVOTAL_DIR}/*.pivotal)
PIVOTAL_BASENAME=$(basename ${PIVOTAL_FILE})

remote_folder=${1:-untested}

bundle
bundle exec rake opsmgr:bucket:put[${remote_folder},${PIVOTAL_FILE}.md5]
bundle exec rake opsmgr:bucket:put[${remote_folder},${PIVOTAL_FILE}.yml]
bundle exec rake opsmgr:bucket:put[${remote_folder},${PIVOTAL_FILE}]

echo "${remote_folder}/${PIVOTAL_BASENAME}.md5" > "${PIVOTAL_DIR}/cf-pivotal.md5.blobkey"
echo "${remote_folder}/${PIVOTAL_BASENAME}.yml" > "${PIVOTAL_DIR}/cf-pivotal.yml.blobkey"
echo "${remote_folder}/${PIVOTAL_BASENAME}" > "${PIVOTAL_DIR}/cf-pivotal.blobkey"
