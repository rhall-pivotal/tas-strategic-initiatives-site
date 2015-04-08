#!/bin/bash

set -ex

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd ../ && pwd )"

bundle

mkdir -p "${SCRIPTS_DIR}/../../opsmgr-installation/"
INSTALLATION_DIR="$(cd ${SCRIPTS_DIR}/../../opsmgr-installation/ && pwd)"
INSTALLATION_ZIP="${INSTALLATION_DIR}/installation.zip"

curl -k $(bundle exec rake --trace opsmgr:info:export_installation_url[${RELENG_ENV}]) -o ${INSTALLATION_ZIP} &&
  unzip -l ${INSTALLATION_ZIP} && # seems to always fail on OSX when the zip file exceeds 4GB
   ruby -r 'digest/md5' -e "puts Digest::MD5.file('${INSTALLATION_ZIP}').hexdigest" > ${INSTALLATION_ZIP}.md5
