#!/bin/bash

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

INSTALLATION_ZIP="$(echo ${SCRIPTS_DIR}/../../opsmgr-installation/*.zip)"
INSTALLATION_ZIP_MD5="$(echo ${SCRIPTS_DIR}/../../opsmgr-installation/*.zip.md5)"

EXPECTED_MD5_VALUE=$(cat "${INSTALLATION_ZIP_MD5}")

ACTUAL_MD5_VALUE=$(ruby -r 'digest/md5' -e "puts Digest::MD5.file('${INSTALLATION_ZIP}').hexdigest")
echo "MD5 of new installation.zip is ${ACTUAL_MD5_VALUE}"
if [[ "${ACTUAL_MD5_VALUE}" != "${EXPECTED_MD5_VALUE}" ]]; then
  echo "installation.zip does not match md5" >&2
  exit 1
fi

INSTALLATION_URL=$(bundle exec rake --trace opsmgr:info:import_installation_url[${RELENG_ENV}])
OPSMGR_PASSWORD=$(bundle exec rake --trace opsmgr:info:password[${RELENG_ENV}])

echo "Uploading installation.zip to Ops Manager."
curl -k ${INSTALLATION_URL} -X POST -F "password=${OPSMGR_PASSWORD}" -F "installation[file]=@${INSTALLATION_ZIP}"
