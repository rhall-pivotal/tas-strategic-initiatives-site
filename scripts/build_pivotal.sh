#!/bin/bash

set -e

P_RUNTIME_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && cd .. && pwd )"

bundle install
bundle exec vara-build-metadata --product-dir="${P_RUNTIME_DIR}"
bundle exec vara-download-artifacts --product-metadata="${P_RUNTIME_DIR}/metadata/cf.yml"
bundle exec vara-build-pivotal --product-metadata="${P_RUNTIME_DIR}/metadata/cf.yml" --rc="-build${BUILD_NUMBER:--local}"
